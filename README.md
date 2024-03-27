______    _     _             _____                            _
|  ___|  | |   | |           |  ___|                          | |
| |_ ___ | | __| | ___ _ __  | |__ _ __   ___ _ __ _   _ _ __ | |_ ___  _ __
|  _/ _ \| |/ _` |/ _ \ '__| |  __| '_ \ / __| '__| | | | '_ \| __/ _ \| '__|
| || (_) | | (_| |  __/ |    | |__| | | | (__| |  | |_| | |_) | || (_) | |
\_| \___/|_|\__,_|\___|_|    \____/_| |_|\___|_|   \__, | .__/ \__\___/|_| 
V2 by Robin K.                                      |__/|_| 


Memory-Efficient and parallel FolderEncryptor which can encrypt small and large folders.
Encryption is done using ChaCha20/HMAC. 

In addition to the encryption it also generates a table with the old Filename, 
the new randomly generated fileName, the used nonce, the used salt and the 
Message Authentication Code and encrypts it with ChaCha20-Poly1305.

If the Poly1305 check fails, the program stops to prevent destroying the files 
with the wrong password.

Key Derivation is done using PKDF2.

In the future i will probably add:
-Menu for tweaking the Iteration-count for PKDF2
-Obfuscation of folder names (currently only filenames are obfuscated)
-Maybe a second encryptor for encrypting small files which encrypts them in place 
  (way faster then streaming through small files, which is the current bottleneck)
