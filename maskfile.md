## build

```sh
rm -Rf build
mkdir build
cp -R images/appicon.png build
wails build --platform "darwin/universal"
```
