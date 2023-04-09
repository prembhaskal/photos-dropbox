# Photos to Dropbox
app to transfer my pics from photos to dropbox to avoid reaching limit.


## Notes are stored here
 https://docs.google.com/document/d/1kGZajTi2gXiAvMPcXR9jcpLnHh1fm0a9jALVKAKobe0/edit?usp=sharing

### running in local
```
$ npm install http-server -g
$ cd ishancraft_handmade
$ http-server
```
  
open http://localhost:8080 in chrome/firefox.  

NOTE: The port 19080 exists in code because of my current testing environment in window which is port-fowarding from virtualbox.


### current usage
#### google photos
  - open http://localhost:19080
  - open inspect window and open console to keep monitoring logs for any error
  - and then click on try request, select your google account
  - it will redirect page to localhost with access_token now in URL
  - the table will show the 1st 10 images in small format.
  - click on each image to open it in new window in bigger size.
  - 'Next Data' button does nothing yet.
#### dropbox pics
  - open http://localhost:19080/dropbox.html
  - click on the button 'List Dropbox Folders' 3 times to see the thumbnail and picture of ABD.
  - click on either picture to trigger a download.
  - all interesting things are available in console here, it still under development.