"""全局配置：通过 Pydantic Settings 从 .env 文件加载所有环境变量。"""

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    # 服务端口
    PORT: int = 5557

    # 智谱 GLM API
    GLM_API_KEY: str = ""
    GLM_MODEL_NAME: str = "glm-4.7-flash"
    GLM_API_BASE_URL: str = "https://open.bigmodel.cn/api/paas/v4"
    GLM_EMBEDDING_MODEL: str = "embedding-2"

    # Milvus 向量数据库
    MILVUS_HOST: str = "localhost"
    MILVUS_PORT: int = 19530
    MILVUS_COLLECTION: str = "algorithm_knowledge"

    # Go 主后端（AI 后端调用其 API 获取题目/用户数据）
    GO_BACKEND_URL: str = "http://127.0.0.1:8080"

    # JWT 签名密钥（必须与 Go 后端 middleware/jwt/jwt.go 中的 SignKey 一致）
    JWT_SIGN_KEY: str = "xwq200505123337"

    # CORS
    CORS_ORIGINS: list[str] = ["*"]

    model_config = {"env_file": ".env", "env_file_encoding": "utf-8", "extra": "ignore"}


settings = Settings()
