# fhana

A lightweight and simple **ransomware** implementation written in Go language, for security testing or educational purposes. Check if the security software or system logger can detect this code's behavior properly.

In `.env`, there is an environment variable named `TARGET_ABSOLUTE_DIRPATH` which specifies the target directory of `fhana`. It will recursively explore all existing files under the directory, conducting rapid **AES-GCM encryption** with multiple **goroutines**. The AES key file will be automatically generated as `key.txt`, storing 128-bit hexadecimal values.

```plaintext
TARGET_ABSOLUTE_DIRPATH="C:/Temp/myDocuments"
```

## Summary

- **Recursive File Encryption**: Traverses all subdirectories and files within the specified target directory.
- **AES-GCM Encryption**: Utilizes AES-GCM for secure encryption of file contents.
- **Concurrent Processing**: Implements a worker pool with multiple goroutines to speed up the encryption process.
- **Key Management**: Automatically generates and stores a 128-bit AES key in `key.txt`.

## How It Works

1. **Environment Setup**:
   - Specify the target directory in the `.env` file.
   - Example:
     ```plaintext
     TARGET_ABSOLUTE_DIRPATH="C:/Temp/myDocuments"
     ```

2. **Encryption Process**:
   - The program reads the `TARGET_ABSOLUTE_DIRPATH` from the environment variable.
   - Generates a 128-bit AES key and saves it in `key.txt`.
   - Uses `filepath.Walk` to recursively list all files in the target directory.
   - Spawns multiple worker goroutines (default is 10) to encrypt files concurrently.
   - Each file is encrypted using AES-GCM, and the original file is securely deleted after encryption.

## Important Notes

- **WARNING**: Be extremely careful not to encrypt important directories (`TARGET_ABSOLUTE_DIRPATH`). Careless execution of this code can result in **irrevocable and destructive** consequences.
- **Educational Use Only**: This code is intended for educational purposes or security testing. Misuse of this code for malicious purposes is illegal and unethical.

## Usage

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/yourusername/fhana.git
   cd fhana
   ```

2. **Initialize Go Modules**:
   ```sh
   go mod tidy
   ```

3. **Set Up Environment Variables**:
   - Create a `.env` file in the project root directory with the following content:
     ```plaintext
     TARGET_ABSOLUTE_DIRPATH="C:/Temp/myDocuments"
     ```

4. **Run the Program**:
   ```sh
   go run main.go
   ```

5. **Decrypt Files**:
   - To decrypt files, use the `-decrypt` flag:
     ```sh
     go run main.go -decrypt
     ```
   - **IMPORTANT**: You have to possess the key file that was used to **the encryption procedure**, and it has to be placed as **`<PROJECT_DIR>/key.txt`**. **Using wrong key files for the decryption may result irrevocable consequences**. Also, it's technically impossible to decrypt files without the key file too.

By following these steps, you can test the encryption process and observe how security tools respond to the ransomware behavior. Remember, this tool is for learning and testing purposes only. Use it responsibly.