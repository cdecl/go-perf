import os
import psutil
import socket
import json
import threading
from datetime import datetime
    
def datetime14():
	return datetime.now().strftime('%Y%m%d%H%M%S')

def curPerf():
	ts = datetime14()
		#datetime.isoformat()

	perf = {
		#'@timestamp': ts,
		'timestamp': ts,
		'hostname': socket.gethostname(),
		'processor': psutil.cpu_percent(interval=1),
		'memory': psutil.virtual_memory().percent,
		'swap': psutil.swap_memory().percent,
		'disk': psutil.disk_usage('/').percent,
		'ip': socket.gethostbyname(socket.getfqdn())
	}

	try:
		perf["loadavg"] = psutil.getloadavg()[1]
	except Exception:
		pass

	return json.dumps(perf)

def main():
	perf = curPerf()
	fname = "perf-{}.log".format(datetime14()[0:8])

	print(perf)

if __name__ == "__main__":
	main()
