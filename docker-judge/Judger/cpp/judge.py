from utils import create_file,extract_log_info,generateRandomPathString,generateRandomDockerName
from pathlib import Path
import subprocess
import os
import asyncio
import json
DOCKER_IMAGE = "sandbox_universal"  # Docker镜像名称
TIME_PATH = '/usr/bin/time'
CODE_PATH = "/tmp"
with open("./config/file_suffix.json", 'r', encoding='utf-8') as file:
    file_suffix = json.load(file) 
    
async def C_cpp(code, language, JudgeCase, time_limit=20, memory_limit=512):
    temp_dir = generateRandomPathString()
    os.makedirs(temp_dir)
    code_file = temp_dir / f"Main.{file_suffix[language]}"
    exec_path = temp_dir/ "Main"
    err_path = temp_dir / "err"
    input_file = temp_dir / "in"
    output_file = temp_dir / "out"
    log_file = temp_dir / "log"
    # print("输入数据: ", input_file)
    # print("输出数据: " , output_file)
    # print("错误报告: ",err_path)
    # print("性能报告: ",log_file)
    await asyncio.gather(
        create_file(code_file, code),
        create_file(err_path),
        create_file(input_file, JudgeCase),
        create_file(output_file),
        create_file(log_file)
    )
    env_cmd = "export LD_LIBRARY_PATH=/usr/lib/x86_64-linux-gnu:$LD_LIBRARY_PATH"
    compile_cmd = {
        "c": f"gcc -o {exec_path} {code_file} 2> {err_path}  ",
        "cpp": f"g++ -std=c++2a -o {exec_path} {code_file}  -lfolly -lglog -lgflags -lpthread -ldl -ldouble-conversion -liberty -lfmt 2> {err_path} "
    }.get(language)
    run_cmd = f"{TIME_PATH} -v sh -c \"{exec_path} < {input_file} > {output_file} 2>> {err_path}\" 2> {log_file}"
    all_cmd = f"{env_cmd} && {compile_cmd} && {run_cmd}"
    mount_cmd = f" -v {temp_dir}:{temp_dir} -v {TIME_PATH}:{TIME_PATH}"
    command = (
            f"sudo docker run --rm -m {memory_limit}M  -i {mount_cmd} --name {generateRandomDockerName()} {DOCKER_IMAGE} "
            f"sh -c '{all_cmd} '")
    try:
        subprocess.run(command,timeout = time_limit, shell=True, check=True)
    except subprocess.TimeoutExpired:
        return {'info': [], 'err': "time limit exceed ", 'msg' : "超时"}
    except subprocess.CalledProcessError as e:
        with open(err_path, "r") as f:
            return {'info': [], 'err': f.read(),'msg' : '发生错误'}
    res = {}
    for file, key in [(input_file, 'input'), (output_file, 'output'), (err_path, 'error')]:
        with open(file, 'r') as f:
            res[key] = f.read()
    
    res['time'], res['memory'] = extract_log_info(log_file)
    res['error'] = 0 if all(k in res for k in ['time', 'memory']) else "Internal Server Error"

    os.system(f"sudo rm -rf {temp_dir}") 
    return {'info': res, 'err': '', 'msg': ''} 

