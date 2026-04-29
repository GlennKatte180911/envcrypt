# envcrypt

> A simple utility to encrypt and manage `.env` files for teams using symmetric keys.

---

## Installation

```bash
go install github.com/yourusername/envcrypt@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envcrypt.git && cd envcrypt && go build ./...
```

---

## Usage

**Encrypt a `.env` file:**

```bash
envcrypt encrypt --key mysecretkey --in .env --out .env.enc
```

**Decrypt a `.env` file:**

```bash
envcrypt decrypt --key mysecretkey --in .env.enc --out .env
```

**Generate a secure key:**

```bash
envcrypt keygen
```

Share the encrypted `.env.enc` file safely in version control and distribute the key to teammates via a secrets manager or secure channel.

---

## How It Works

`envcrypt` uses AES-256-GCM symmetric encryption to protect your `.env` files. Each encryption operation produces a unique nonce, ensuring the same input never produces the same output twice.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername