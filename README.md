Memory-Efficient and parallel FolderEncryptor which can encrypt small and large folders.<br>
Encryption is done using ChaCha20/HMAC. <br>

In addition to the encryption it also generates a table with the old filename, <br>
the new randomly generated filename, the used nonce, the used salt and the <br>
Message Authentication Code. The table also gets encrypted for obfuscation <br>
with ChaCha20-Poly1305.<br>
<br>
If the Poly1305 check fails, the program stops to prevent destroying the files <br>
with the wrong password.<br>
<br>
Key Derivation is done using PKDF2.<br>
<br>
In the future i will probably add:<br>
-Menu for tweaking the Iteration-count for PKDF2<br>
-Obfuscation of folder names (currently only filenames are obfuscated)<br>
-Maybe a second encryptor for encrypting small files which encrypts them in place <br>
  (way faster then streaming through small files, which is the current bottleneck)<br>
