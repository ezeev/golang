
echo "Regular gauge"
echo "gauge1:10|g" | nc -u -w0 localhost 10001
echo "gauge1:10|g" | nc -u -w0 localhost 10001
echo "gauge1:10|g" | nc -u -w0 localhost 10001

echo "Incrementing Gauge"
echo "gauge2:+10|g" | nc -u -w0 localhost 10001
echo "gauge2:+10|g" | nc -u -w0 localhost 10001
echo "gauge2:+10|g" | nc -u -w0 localhost 10001

echo "Decrementing Gauge"
echo "gauge2:-10|g" | nc -u -w0 localhost 10001
echo "gauge2:-10|g" | nc -u -w0 localhost 10001
echo "gauge2:-10|g" | nc -u -w0 localhost 10001

echo "Counter"
echo "counter1:10|c" | nc -u -w0 localhost 10001
echo "counter1:10|c" | nc -u -w0 localhost 10001
echo "counter1:10|c" | nc -u -w0 localhost 10001

echo "Timer"
echo "timer1:20|ms" | nc -u -w0 localhost 10001
echo "timer1:10|ms" | nc -u -w0 localhost 10001
echo "timer1:30|ms" | nc -u -w0 localhost 10001
echo "timer1:42|ms" | nc -u -w0 localhost 10001
echo "timer1:2|ms" | nc -u -w0 localhost 10001
echo "timer1:60|ms" | nc -u -w0 localhost 10001
echo "timer1:21|ms" | nc -u -w0 localhost 10001
