import re
import random
import string
import time
from pathlib import Path
CODE_PATH = "/tmp" 
async def create_file(path, content=None):
    with open(path, "w") as file:
        if content is not None:
            file.write(content)


def extract_log_info(log_file_path):
    with open(log_file_path, 'r') as file:
        log_content = file.read()
    elapsed_time_pattern = re.compile(r'Elapsed \(wall clock\) time \(h:mm:ss or m:ss\): ([0-9:]+\.[0-9]+)')
    max_resident_set_size_pattern = re.compile(r'Maximum resident set size \(kbytes\): ([0-9]+)')
    elapsed_time_match = elapsed_time_pattern.search(log_content)
    max_resident_set_size_match = max_resident_set_size_pattern.search(log_content)
    elapsed_time = None
    if elapsed_time_match:
        elapsed_time_str = elapsed_time_match.group(1)
        parts = elapsed_time_str.split(':')
        if len(parts) == 2:  # 格式为 m:ss 或 h:mm:ss
            minutes, seconds = map(float, parts)
            elapsed_time = minutes * 60 + seconds
        elif len(parts) == 3:  # 格式为 h:mm:ss
            hours, minutes, seconds = map(float, parts)
            elapsed_time = hours * 3600 + minutes * 60 + seconds
        else:
            elapsed_time = float(elapsed_time_str)  # 如果只有一个部分，则直接是秒数

    max_resident_set_size = int(max_resident_set_size_match.group(1)) if max_resident_set_size_match else None
    return elapsed_time, max_resident_set_size


def generateRandomPathString():
    characters = string.ascii_letters + string.digits 
    timestamp_part = str(int(time.time()))[4:]
    random_part = ''.join(random.choices(characters, k=5))
    final_path = Path(CODE_PATH) / (timestamp_part + random_part)
    return final_path


def generateRandomDockerName():
    characters = string.ascii_letters + string.digits
    timestamp_part = str(int(time.time()))[6:]
    random_part = ''.join(random.choices(characters, k=6))
    return random_part +  timestamp_part 
    