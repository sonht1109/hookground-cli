# Hookground CLI
This is a CLI tool for [hookground](https://hookground.sonht.io.vn) which allows you to forward webhook to your local server

>This project is only for study purpose.

### Installation
1. Go to https://hookground.sonht.io.vn and try receiving some webhooks
2. Visit webhook detail
3. Run CLI
```bash
curl -sSL https://sonht1109.github.io/hookground-cli/install.sh | bash -s -- --t <YOUR_WEBHOOK_KEY> --h <LOCAL_TARGET_URL>
```
4. Your target host (eg: http://localhost:3000/path) will receive webhook

### Note
Make sure that your local server is running and also allows POST method