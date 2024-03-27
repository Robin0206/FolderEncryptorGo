Memory-Efficient and parallel FolderEncryptor which can encrypt small and large folders.
Encryption is done using ChaCha20/HMAC. 

In addition to the encryption it also generates a table with the old filename, 
the new randomly generated filename, the used nonce, the used salt and the 
Message Authentication Code. The table also gets encrypted for obfuscation 
with ChaCha20-Poly1305.

If the Poly1305 check fails, the program stops to prevent destroying the files 
with the wrong password.

Key Derivation is done using PKDF2.

In the future i will probably add:\n
-Menu for tweaking the Iteration-count for PKDF2
-Obfuscation of folder names (currently only filenames are obfuscated)
-Maybe a second encryptor for encrypting small files which encrypts them in place 
  (way faster then streaming through small files, which is the current bottleneck)
