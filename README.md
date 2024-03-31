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

Usage on Linux:<br>
<br>
There is a binary in the out folder.<br>
Additionally you still can compile your own one using go build.<br>

<br>
Add the following line to the end of your .bashrc<br>
alias fenc='/path/to/the/binary/NameOfTheCompiledBinary ${PWD}'<br>
<br>
cd into the folder and type fenc<br>
<br>
Usage on Windows is untested but the program takes the path to the folder as its first argument.<br>
In theory it should just work if you compile your own binary using go build(on a windows machine).<br>
After that add the path to the folder thats containing your exe to the path variable.<br>
