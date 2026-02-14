import os
import glob
from dotenv import load_dotenv
from pymilvus import (
    connections,
    Collection,
    FieldSchema,
    CollectionSchema,
    DataType,
    utility
)
from zhipuai import ZhipuAI
from loguru import logger
from tqdm import tqdm
import fitz  # PyMuPDF

# 配置日志
logger.add("log/vectorization.log", rotation="10 MB", retention="10 days", level="INFO")

# 加载环境变量
load_dotenv()

# Milvus连接配置
MILVUS_CONFIG = {
    "host": os.getenv("MILVUS_HOST", "localhost"),
    "port": os.getenv("MILVUS_PORT", "19530"),
    "user": os.getenv("MILVUS_USER", ""),
    "password": os.getenv("MILVUS_PASSWORD", ""),
}

# ZhipuAI配置
ZHIPU_API_KEY = os.getenv("ZHIPU_API_KEY")
client = ZhipuAI(api_key=ZHIPU_API_KEY)

# 集合名称
COLLECTION_NAME = "algorithm_knowledge"

# Embedding维度（ZhipuAI的embedding-2模型输出1024维向量）
EMBEDDING_DIM = 1024


def connect_milvus():
    """连接到Milvus数据库"""
    try:
        # 阿里云Milvus使用host:port的方式连接
        logger.info(f"尝试连接Milvus: {MILVUS_CONFIG['host']}:{MILVUS_CONFIG['port']}")
        
        connections.connect(
            alias="default",
            host=MILVUS_CONFIG["host"],
            port=MILVUS_CONFIG["port"],
            user=MILVUS_CONFIG["user"],
            password=MILVUS_CONFIG["password"],
        )
        logger.info(f"成功连接到Milvus数据库")
        return True
    except Exception as e:
        logger.error(f"连接Milvus失败: {e}")
        return False


def create_collection():
    """创建Milvus集合"""
    try:
        # 如果集合已存在，先删除
        if utility.has_collection(COLLECTION_NAME):
            logger.info(f"集合 {COLLECTION_NAME} 已存在，将删除后重新创建")
            _ = utility.drop_collection(COLLECTION_NAME)
        
        # 定义字段
        fields = [
            FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
            FieldSchema(name="file_name", dtype=DataType.VARCHAR, max_length=500),
            FieldSchema(name="content", dtype=DataType.VARCHAR, max_length=65535),
            FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=EMBEDDING_DIM)
        ]
        
        # 创建集合schema
        schema = CollectionSchema(fields=fields, description="程序设计竞赛中的算法知识")
        
        # 创建集合
        collection = Collection(name=COLLECTION_NAME, schema=schema)
        logger.info(f"成功创建集合: {COLLECTION_NAME}")
        
        # 创建索引
        index_params = {
            "metric_type": "COSINE",  # 使用余弦相似度
            "index_type": "IVF_FLAT",
            "params": {"nlist": 128}
        }
        _ = collection.create_index(field_name="embedding", index_params=index_params)
        logger.info("成功创建索引")
        
        return collection
    except Exception as e:
        logger.error(f"创建集合失败: {e}")
        return None


def get_embedding(text):
    """使用ZhipuAI获取文本的向量表示"""
    try:
        response = client.embeddings.create(
            model="embedding-2",
            input=text
        )
        embedding = response.data[0].embedding
        return embedding
    except Exception as e:
        logger.error(f"获取embedding失败: {e}")
        return None


def read_files(directory):
    """递归读取目录下的所有txt和pdf文件"""
    data = []
    
    # 使用os.walk递归遍历目录
    for root, dirs, files in os.walk(directory):
        txt_files = [f for f in files if f.endswith('.txt')]
        pdf_files = [f for f in files if f.endswith('.pdf')]
        logger.info(f"在目录 {root} 中找到 {len(txt_files)} 个txt文件和 {len(pdf_files)} 个pdf文件")
        
        # 处理txt文件
        for file_name in txt_files:
            file_path = os.path.join(root, file_name)
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    content = f.read().strip()
                    # 过滤掉空文件
                    if content:
                        # 使用相对路径作为文件标识
                        relative_path = os.path.relpath(file_path, directory)
                        data.append({
                            "file_name": relative_path,
                            "content": content
                        })
            except Exception as e:
                logger.error(f"读取文件 {file_path} 失败: {e}")
        
        # 处理pdf文件
        for file_name in pdf_files:
            file_path = os.path.join(root, file_name)
            try:
                doc = fitz.open(file_path)
                content = ""
                for page in doc:
                    content += page.get_text("text")
                doc.close()
                
                # 过滤掉空文件
                if content.strip():
                    # 使用相对路径作为文件标识
                    relative_path = os.path.relpath(file_path, directory)
                    data.append({
                        "file_name": relative_path,
                        "content": content.strip()
                    })
            except Exception as e:
                logger.error(f"读取PDF文件 {file_path} 失败: {e}")
    
    logger.info(f"成功读取 {len(data)} 个有效txt/pdf文件")
    return data


def vectorize_and_insert(collection, data, batch_size=10):
    """将新闻数据向量化并批量插入Milvus"""
    total = len(data)
    logger.info(f"开始向量化和插入 {total} 条新闻数据")
    
    file_names = []
    contents = []
    embeddings = []
    
    # 使用进度条
    with tqdm(total=total, desc="向量化进度") as pbar:
        for idx, item in enumerate(data):
            # 获取embedding
            embedding = get_embedding(item["content"])
            
            if embedding:
                file_names.append(item["file_name"])
                contents.append(item["content"])
                embeddings.append(embedding)
                
                # 批量插入
                if len(file_names) >= batch_size:
                    try:
                        collection.insert([file_names, contents, embeddings])
                        collection.flush()
                        logger.info(f"成功插入 {len(file_names)} 条数据")
                        file_names = []
                        contents = []
                        embeddings = []
                    except Exception as e:
                        logger.error(f"插入数据失败: {e}")
            else:
                logger.warning(f"跳过文件 {item['file_name']}")
            
            pbar.update(1)
    
    # 插入剩余数据
    if file_names:
        try:
            collection.insert([file_names, contents, embeddings])
            collection.flush()
            logger.info(f"成功插入剩余 {len(file_names)} 条数据")
        except Exception as e:
            logger.error(f"插入剩余数据失败: {e}")
    
    # 加载集合到内存以支持搜索
    collection.load()
    logger.info("集合已加载到内存，可以进行搜索")
    
    # 获取集合统计信息
    num_entities = collection.num_entities
    logger.info(f"集合中共有 {num_entities} 条数据")
    
    return num_entities


def main():
    """主函数"""
    logger.info("=" * 50)
    logger.info("开始执行数据向量化入库流程")
    logger.info("=" * 50)
    
    # 1. 连接Milvus
    if not connect_milvus():
        logger.error("无法连接到Milvus，程序退出")
        return
    
    # 2. 创建集合
    collection = create_collection()
    if not collection:
        logger.error("无法创建集合，程序退出")
        return
    
    # 3. 读取文件
    directory = "DATA"
    data = read_files(directory)
    
    if not data:
        logger.error("没有读取到有效的数据，程序退出")
        return
    
    # 4. 向量化并插入数据
    num_inserted = vectorize_and_insert(collection, data)
    
    logger.info("=" * 50)
    logger.info(f"数据向量化入库完成！共插入 {num_inserted} 条数据")
    logger.info("=" * 50)


if __name__ == "__main__":
    main()
