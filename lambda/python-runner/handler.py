import os
import sys
import time
import resource
import subprocess
import tempfile

MAX_TIMEOUT_S = 8          # hard ceiling regardless of requested timeout
CPU_SECONDS = 5            # RLIMIT_CPU (SIGXCPU past this)
MEM_BYTES = 256 * 1024 * 1024
MAX_PROCS = 64             # blocks fork bombs
MAX_FILE_BYTES = 8 * 1024 * 1024
OUTPUT_CAP = 64 * 1024     # bytes per stream returned

PY = sys.executable        # the runtime's python, e.g. /var/lang/bin/python3


def _apply_limits():
    # Runs in the child, after fork, before exec.
    resource.setrlimit(resource.RLIMIT_CPU, (CPU_SECONDS, CPU_SECONDS + 1))
    resource.setrlimit(resource.RLIMIT_AS, (MEM_BYTES, MEM_BYTES))
    resource.setrlimit(resource.RLIMIT_NPROC, (MAX_PROCS, MAX_PROCS))
    resource.setrlimit(resource.RLIMIT_FSIZE, (MAX_FILE_BYTES, MAX_FILE_BYTES))
    os.setsid()  # own process group so a timeout can kill the whole tree


def _clip(b):
    return (b or b"")[:OUTPUT_CAP].decode("utf-8", errors="replace")


def run_once(source: str, stdin: str, timeout_s: float) -> dict:
    started = time.time()
    with tempfile.TemporaryDirectory() as workdir:
        main_path = os.path.join(workdir, "main.py")
        with open(main_path, "w", encoding="utf-8") as f:
            f.write(source)

        # Minimal environment — no AWS_* creds, no inherited secrets.
        env = {
            "PATH": os.path.dirname(PY) + ":/usr/bin:/bin",
            "HOME": workdir,
            "TMPDIR": workdir,
            "LANG": "C.UTF-8",
            "PYTHONDONTWRITEBYTECODE": "1",
        }
        try:
            proc = subprocess.run(
                [PY, "-I", main_path],          # -I: isolated mode
                input=stdin.encode("utf-8"),
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                cwd=workdir,
                env=env,
                timeout=timeout_s,
                preexec_fn=_apply_limits,
                check=False,
            )
            return {
                "stdout": _clip(proc.stdout),
                "stderr": _clip(proc.stderr),
                "exitCode": proc.returncode,
                "timedOut": False,
                "durationMs": int((time.time() - started) * 1000),
            }
        except subprocess.TimeoutExpired as e:
            return {
                "stdout": _clip(e.stdout),
                "stderr": _clip(e.stderr) + "\n[Killed: time limit exceeded]",
                "exitCode": -1,
                "timedOut": True,
                "durationMs": int((time.time() - started) * 1000),
            }


def handler(event, context):
    mode = event.get("mode", "run")
    source = event.get("source", "") or ""
    timeout_s = min(float(event.get("timeoutMs", 8000)) / 1000.0, MAX_TIMEOUT_S)

    if mode == "grade":
        results = [run_once(source, t.get("input", "") or "", timeout_s)
                   for t in event.get("tests", [])]
        return {"results": results}

    return run_once(source, event.get("stdin", "") or "", timeout_s)
