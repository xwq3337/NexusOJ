from utils import create_file,extract_log_info,generateRandomPathString,generateRandomDockerName
from pathlib import Path
import subprocess
import os
import asyncio
import json
# DOCKER_IMAGE = "sandbox_go"
DOCKER_IMAGE = "sandbox_universal"
TIME_PATH = '/usr/bin/time'
CODE_PATH = "/tmp"
async def Go(code, language, JudgeCase, time_limit=200, memory_limit=256):
    temp_dir = generateRandomPathString()
    os.makedirs(temp_dir)
    code_file = temp_dir / "Main.go"
    exec_path = temp_dir / "Main"
    err_path = temp_dir / "err"
    input_file = temp_dir / "in"
    output_file = temp_dir / "out"
    log_file = temp_dir / "log"
    # print("代码文件: " , code_file)
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
    compile_cmd = f"go build -o {exec_path} {code_file} 2> {err_path}"
    mount_cmd = f" -v {temp_dir}:{temp_dir} -v {TIME_PATH}:{TIME_PATH}"
    run_cmd = f"{TIME_PATH} -o {log_file} -v {exec_path} < {input_file} > {output_file} 2> {err_path}"
    command = (
            f"sudo docker run --rm -m {memory_limit}M  -i {mount_cmd} --name {generateRandomDockerName()} {DOCKER_IMAGE}"
            f" sh -c '{compile_cmd} && {run_cmd}'")
    try:
        subprocess.run(command, shell=True, check=True)
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
    result = {'info': res, 'err': '', 'msg': ''}
    # os.system(f"sudo rm -rf {temp_dir}") 
    return result 
