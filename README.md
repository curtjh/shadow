# Microsoft RDP Shadow Go implimentation
Microsoft RDP shadow app in Go.

Microsoft provides a way to shadow a user's session out of the box within mstsc.exe using the /shadow swtich. However, before running that, you'll need to know the session number of the remote user, and to get that you'll need to run another command first (qwinsta /server:[remoteComputer]). This combines those 2 actions (getting remote session id, then attempting a shadow) into one exe so if you use shadowing, this will make your life easier.

Requirements:

All that's required on for this to work is that you have permissions to run 'qwinsta.exe' and 'rwinsta.exe' on remote computers.
Example to test from your command line:
```
qwinsta /server:[computer]
```
if that produces result without error, you can use this provided RDP Shadowing is enabled on the target computer

Examples:

1. Shadow a computer with no control and no consent:
   ```
   shadow.exe -v [computername] 
   ```

2. Shadow a computer with control and consent:

   ```
   shadow.exe -v [computername] -consent -control
   ```
   
3. Shadow a computer with control but no consent:
   ```
   shadow.exe -v [computername] -control
   ```  
4. Shadow a computer with consent but no control:
   ```
   shadow.exe -v [computername] -consent
   ```
5. List connected users: (use of this switch will not shadow a session, only list users)
   ```
   shadow.exe -v [computername] -listUsers 
   ```
6. Disconnect a session: (use of this switch will not shadow a session, only attempt to disconnect passed session number)
   ```
   shadow.exe -v [computername] -dissconnect [sessionNumber] 
   ```
   
7. Prompt for credentials: (use this if your account can't shadow, but you have access to an account that can)
   ```
   shadow.exe -v [computername] -credPrompt 
   ```
