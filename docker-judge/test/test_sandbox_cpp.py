import os

os.system("clear")
os.system("rm -f err log out main")
PATH = "/home/ubuntu/judge/test"
TIME_PATH = "/usr/bin/time"
mount_cmd = f"  -v {PATH}:{PATH}"
DOCKER_IMAGE = "sandbox_cpp_plus"
env_cmd = "sh -c 'export LD_LIBRARY_PATH=/usr/lib/x86_64-linux-gnu:$LD_LIBRARY_PATH'"
compile_cmd = f"g++ -o {PATH}/main 2> err {PATH}/1.cpp -I/usr/include -I/usr/local/include -L/usr/lib -L/usr/local/lib -lfolly -lglog -lgflags -lpthread -ldl -ldouble-conversion -liberty -lfmt "
ln_cmd = "ln -s /usr/lib/x86_64-linux-gnu/libfmt.so.8.1.1 /usr/lib/x86_64-linux-gnu/libfmt.so.9"
run_cmd = f" {PATH}/main < {PATH}/in > {PATH}/out 2> {PATH}/log"

print("编译命令： ",  compile_cmd)
print("运行命令： ",  run_cmd)
print("总命令： ",f"sudo docker run --rm -i {mount_cmd} {DOCKER_IMAGE} {env_cmd} && {compile_cmd} && {run_cmd} ")
os.system(f"sudo docker run --rm -i {mount_cmd} {DOCKER_IMAGE} {env_cmd} && {compile_cmd} && {run_cmd} ")

"""
sudo docker run --rm -i \
  -v /home/ubuntu/judge/test:/home/ubuntu/judge/test \
  sandbox_cpp_plus sh -c 'export LD_LIBRARY_PATH=/usr/lib/x86_64-linux-gnu:$LD_LIBRARY_PATH && \
  g++ -o /home/ubuntu/judge/test/main /home/ubuntu/judge/test/1.cpp -I/usr/include -I/usr/local/include -L/usr/lib -L/usr/local/lib -lfolly -lglog -lgflags -lpthread -ldl -ldouble-conversion -liberty -lfmt 2> /home/ubuntu/judge/test/err && \
  /home/ubuntu/judge/test/main < /home/ubuntu/judge/test/in > /home/ubuntu/judge/test/out 2> /home/ubuntu/judge/test/log'
"""

