# Home

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
    - [Mac OS](#mac-os-installation)
    - [Windows](#windows-installation)
- [Setup](#first-time-setup)
- [Usage](#usage)
    - [View Data](#view-data)
    - [Send Data](#send-data)
    - [Settings](#settings)
    - [Exit](#exit)

## Introduction

Freeport is a TUI that allows for seamless data transfer between localhosted services. 

# Installation

To install Freeport, go to the [releases](https://github.com/kashsuks/freeport/releases) tab on the [Github Repository](https://github.com/kashsuks/freeport) and install whichever fits your operating system!

**Note: Freeport currently has support for Mac Silicon and Windows 64-Bit, there are more supported versions coming soon.**

### Mac OS Installation

Setting up the application for MacOS is quite easy. Once you have the latest release from the [releases tab](https://github.com/kashsuks/freeport/releases), make sure you copy the *full path* of the installed file.

Go to the terminal and type:
```
chmod +x <path of the install>
```

Once that is done *and there are no errors*, the app should be compiled for your computer. To check, go to the folder in which the file is located and type `ls`. If the app name appears then you are good to go!

You may now type in 
```
./freeport
```
And the application should open up.

Conga rats! You just installed freeport :D

### Windows Installation

The installation process for Windows is quite simple!

If you haven't already, head over to the [releases page](https://github.com/kashsuks/freeport/releases) and install the latest `.exe` release! **Ensure that your system uses a 64-Bit processor**.

If windows deems the app to be unsafe, simple ignore the warnings and install the app.

Once that is done, double click the file and it *should* open up in the terminal.

Conga rats! You just installed freeport :D

<div class="tenor-gif-embed" data-postid="7323300" data-share-method="host" data-aspect-ratio="1.745" data-width="100%"><a href="https://tenor.com/view/congrats-congarats-yay-conga-line-conga-gif-7323300">Congrats Congarats GIF</a>from <a href="https://tenor.com/search/congrats-gifs">Congrats GIFs</a></div> <script type="text/javascript" async src="https://tenor.com/embed.js"></script>

## First time setup

Now that you have installed [Freeport]() there is one tiny setup change that you *might* want to make.

That is the welcome message! At the top of the app you will notice that it says
> Welcome to freeport!
This section is fully customizeable and I encourage you to play around with it since it gives you a chance to figure out the controls!

## Usage

After running the app you will see `View Data`, `Send Data`, `Settings`, `Exit`.

### View Data

View Data allows you to view any HTTP method that is accessible by you!

By default you will see that view batter data is displayed. Go ahead and press enter while selecting it! You will see that the battery percentage of your device shows up along with other data such as time and app name.

As of now, this is the only field for view data, but as the project progresses there are more features to be seen here

### Send Data

The Send Data option is the most important one (in my opinion) since it allows you to create protocols!

> What is a protocol?

Well, a protocol allows for your own applications to send data through the pipeline.

The pipeline works by using an `app_name` which is a unique identifier for your specific application. That along with a `passkey` ensure that only authorized users may send data to the application. 

You can use http request headers to send data and the reciever being able to recieve it.

**TL;DR: Protocols are fancy words for http headers!**

Now, I want *you* to try and make your own protocol and send data.

1. Create a new protocol by going to Send Data -> Press `c` for new protocol.
2. Name the protocol `test` and make the passkey `test`. You can make the description whatever you would like it to be
3. Now that the protocol has been created, try testing the application by using curl. Try the following command:
```
curl -H "X-App-Name: test" -H "X-Passkey: test" http://localhost:6767/test/init
```
If the app returns something similar to
```
{"app_name":"test","message":"Hello, World!","status":"initialized","time":"2025-12-22T20:04:57-05:00"}
```
Then the protocol has beeen created successfully!
4. Try creating your own method by going to Freeport and presing `n`. This creates a new method for you to send a recieve data through.
5. Name the method `test` and customize the description however you want.
6. Go back to the terminal and try the following command:
```
curl -X POST \
  -H "X-App-Name: test" \
  -H "X-Passkey: test" \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello!"}' \
  http://localhost:6767/test/test
```
If the request returns something like
```
{"app_name":"test","message":"Data stored successfully","method":"test","status":"success","timestamp":"2025-12-22T20:08:33-05:00"}
```
Then the data has been sent into the pipeline successfully! It will be stored there until another app retrieves it.

7. Try making a get request to the same data by using the following command!
```
curl -H "X-App-Name: test" -H "X-Passkey: test" http://localhost:6767/test/test
```
If the data returned look like the following then it works!
```
{"app_name":"test","data":{"message":"Hello!"},"method":"test","status":"success","time":"2025-12-22T20:11:01-05:00"}
```
This means that the data has been sent to the intended user and the data is dropped from Freeports memory.

Congrats! You just learnt how to use protocols effectively.

### Settings

The settings tab is quite simple for now as it allows you to change the welcome message upon startup.

### Exit

As the name suggests, selecting this option will make you exit the application.