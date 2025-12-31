import json
import asyncio
import os
from pathlib import Path
from utils import create_file,extract_log_info,generateRandomPathString , generateRandomDockerName
import subprocess

DOCKER_IMAGE  = "sandbox_universal"
TIME_PATH = '/usr/bin/time'
CODE_PATH = "/tmp"
async def JavaScript(code, language, JudgeCase, time_limit=1, memory_limit=64):
    temp_dir = generateRandomPathString()
    os.makedirs(temp_dir)
    code_file = temp_dir / "Main.js"
    exec_path = temp_dir / "Main"
    err_path = temp_dir / "err"
    input_file = temp_dir / "in"
    output_file = temp_dir / "out"
    log_file = temp_dir / "log"
    await asyncio.gather(
        create_file(code_file, code),
        create_file(err_path),
        create_file(input_file, JudgeCase),
        create_file(output_file),
        create_file(log_file)
    )
    mount_cmd = f" -v {temp_dir}:{temp_dir} -v {TIME_PATH}:{TIME_PATH}"
    run_cmd = f"{TIME_PATH} -o {log_file} -v node {code_file} < {input_file} > {output_file} 2> {err_path}"
    docker_command = (
        f"sudo docker run --rm -i {mount_cmd} --name {generateRandomDockerName()} {DOCKER_IMAGE}"
        f" sh -c '{run_cmd}'"
    )
    try:
        proc = subprocess.run(docker_command, shell=True, capture_output=True, text=True)
        if proc.returncode != 0 or Path(err_path).read_text().strip():
            return {'info': [], 'err': Path(err_path).read_text(), 'msg': ''}

        res = {
            'input': Path(input_file).read_text(),
            'output': Path(output_file).read_text(),
            'error': '',
            'time': -1,
            'memory': -1,
        }

        res['time'], res['memory'] = extract_log_info(log_file)
        res['error'] = 0 if all(k in res for k in ['time', 'memory']) else "Internal Server Error"
    except Exception as e:
        return {'info': [], 'err': str(e), 'msg': 'Internal Server Error'}
    finally:
        subprocess.run(f"sudo rm -rf {temp_dir}", shell=True)

    return {'info': res, 'err': '', 'msg': ''}
