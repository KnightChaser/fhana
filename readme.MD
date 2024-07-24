# fhana

A lightweight and simple **ransomware** implementation written in Go language, for security testing or educational purposes. Check if the security software or system logger can detect this code's behavior properly.

## Summary

- **Recursive File Encryption**: Traverses all subdirectories and files within the specified target directory.
- **AES-GCM Encryption**: Utilizes AES-GCM for secure encryption of file contents.
- **Concurrent Processing**: Implements a worker pool with multiple goroutines to speed up the encryption process.
- **Key Management**: Automatically generates and stores a 128-bit AES key in the specified key file.

## How It Works

1. **Command-Line Arguments**:
   - **Target Directory**: The directory to be encrypted or decrypted.
   - **Key File**: The file to save the encryption key or load the decryption key.

2. **Encryption Process**:
   - Generates a 128-bit AES key and saves it in the specified key file.
   - Uses `filepath.Walk` to recursively list all files in the target directory.
   - Spawns multiple worker goroutines (default is 10) to encrypt files concurrently.
   - Each file is encrypted using AES-GCM, and the original file is securely deleted after encryption.

3. **Decryption Process**:
   - Reads the AES key from the specified key file.
   - Uses `filepath.Walk` to recursively list all encrypted files in the target directory.
   - Spawns multiple worker goroutines (default is 10) to decrypt files concurrently.
   - Each file is decrypted using AES-GCM, and the encrypted file is securely deleted after decryption.

## Important Notes

- **WARNING**: Be extremely careful not to encrypt important directories. Careless execution of this code can result in **irrevocable and destructive** consequences.
- **Educational Use Only**: This code is intended for educational purposes or security testing. Misuse of this code for malicious purposes is illegal and unethical.

## Usage

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/knightchaser/fhana.git
   cd fhana
   ```

2. **Initialize Go Modules and be ready**:
   ```sh
   go mod tidy
   ```
   - Note: Currently, there are no additional external packages to install. You may skip this for now.

3. **Encryption**:
   - To encrypt files in the target directory and save the key in the specified key file:
     ```powershell
     # Example
     go run main.go --encrypt --target-directory "C:/Temp/myDocument" --key-file "key.txt"
     ```

4. **Decryption**:
   - To decrypt files in the target directory using the specified key file:
     ```powershell
     # Example
     go run main.go --decrypt --target-directory "C:/Temp/myDocument" --key-file "key.txt"
     ```
   - **IMPORTANT**: You must use the same key file that was used during the encryption process. Using the wrong key file for decryption may result in **irrevocable consequences**. It is technically impossible to decrypt files without the correct key file.

By following these steps, you can test the encryption and decryption process and observe how security tools respond to the ransomware behavior. Remember, this tool is for learning and testing purposes only. Use it responsibly.