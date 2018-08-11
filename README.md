# foreman

Service for providing offline update of a raspberry pi via a web interface served by the pi

On receipt of a zip file via its webpage, foreman will:
1. check the hash
2. unzip the file into a temporary folder
3. run the 'install.sh' script

## Build install package

```
make buildpi
```

## Install onto a pi

```
export PI=pi@192.168.1.15
scp install.zip  $PI:/tmp/ && rm install.zip
ssh $PI 'cd /tmp && unzip -o install.zip && rm install.zip && sudo bash install.sh'
```

## Update pi via foreman

For a file called `package.zip` that contains an `install.sh` file:

- Get the hash of the file: `sha256sum package.zip`
- Go to 192.168.1.15:8081
- Upload package.zip and enter the hash
- Output of `install.sh` will be returned when the install has finished
