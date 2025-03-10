import os
import subprocess

def build(platform: str, arch: str, output: str):
    env = os.environ.copy()
    env['GOARCH'] = arch
    env['GOOS'] = platform
    subprocess.run(' '.join([
        'go', 'build', '-ldflags', '"-s -w"', '-o', output
    ]), env=env, check=True)
    print(f'Build {output} success!')

if __name__ == '__main__':
    build('windows', 'amd64', 'bin/minall.exe')
    build('linux', 'amd64', 'bin/minall')