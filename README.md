# youtube-playlist-dl
Downloads youtube playlists in parallel

Install
-------
``` sh
go get github.com/lepidosteus/youtube-dl
```

Examples
-------
``` sh
youtube-playlist-dl "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
```
To specify video quality:
``` sh
youtube-playlist-dl -q hd720 "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
```

You can also request maximum or minimum quality availible:
``` sh
youtube-playlist-dl -q max "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
```
``` sh
youtube-playlist-dl -q min "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
```

Notes
-------
Support for some videos, such as VEVO, is limited.

TODO
-------
 - Specify format (e.g. mp4)
 - Specify number of concurent downloads
