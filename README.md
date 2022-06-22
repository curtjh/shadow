# shadow
Microsoft RDP shadow app in Go.

Examples:

1. Shadow a computer with no control and no consent:

   shadow.exe -v [computername] 

2. Shadow a computer with control and consent:

   shadow.exe -v [computername] -consent -control
   
3. Shadow a computer with control but no consent:

   shadow.exe -v [computername] -control
   
4. Shadow a computer with consent but no control:

   shadow.exe -v [computername] -consent
   
5. List connected users:

   shadow.exe -v [computername] -listUsers (use of this switch will not shadow a session, only list users)
