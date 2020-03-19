# mp3Organizer - Organize your mp3 files
## Q&A
* What does this do?
    * Organize your mp3 files , based on the ID3 tags
    
* What? I didn't get it, explain more!
    * well, you specify a folder, as source and it looks there for 
    
        mp3 files. Also, you can specify destination, which is $HOME/Music
        
        by default if you don't specify, and it moves all those mp3 files to
        
        `$DEST/$Artist/$Album(Optional)(/-)$Title.mp3` where:
        + `$DEST` is destination folder
        + `$Artist` is the name of song artist
        + `$Album` is the name of album (it's optional)
        + `/-` is the option you can choose, weather it creates directory from
            
            last part of destination format, or just rename it.
            
            (`/` is default, so it makes folders unless you explicitly change that)
        + `$Title` is the title of song specified in ID3 Tags
        
* ~~But I don't like how it organizes mp3 files and I love put your template here~~

    ~~Just go to last question, It will answer you!~~

        
* Oh well! How it works then?
    * It works like this: ~~mp3Organizer.exe src dest~~ mp3Organizer.exe
    
        ~~where `src` is your source file and `dest` is destination~~
        
        ~~folder and is optional.~~
        
        You choose source and destination folder in GUI app!
        
* But wait! It has a GUI?
    * Well, I wanted something more visual, and also wanted
    
        to try implementing GUI apps in golang, so here we are.
        
* ~~It's a mix of CLI and GUI? Couldn't it all be GUI?~~
    * ~~Well, I use [fyne](https://github.com/fyne-io/fyne/), and getting source and destination folder~~
    
        ~~in GUI app is not as easy as I thought. It's in my TODO list,~~
        
        ~~and will be fixed if I have time and energy to fix it!~~
        
* I saw your code and it has a strange regex; What does it do?
    * I live in Iran, where *websites* & *telegram channels* are
    
        big source of music downloading and listening!
        
        And it's disgusting that they put their address in ID3 tags!
        
        I mean, who does that, right? So, I had to remove those
        
        addresses from ID3 tags.
        
* Why just window release? No linux?

    * Cause I wrote this program in windows, and compiling it for linux
    
        is not as easy as compiling for windows! I will release a linux binary
        
        as soon as I boot my linux!
        
* I don't like it! It's worthless!

    * I made it for myself first, and shared it as it might help someone!
    
        If it does help someone, then I'm happy, and if not,
        
        I'm using it anyway! So who cares?!
        
        
I hope this *Q&A* answers all your questions!