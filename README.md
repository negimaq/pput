# pput: pre-processing and uploading tool for cloud storage

```bash
cat <<EOF > test
TZ=Asia/Tokyo
ROOT=
USER=
PASSWORD=
DEBUG=false
EOF

sudo docker build -t pput .

sudo docker run --env-file .env -v /path/to/local/directory:/input pput
```
