<div align="center">
  <h1><code>gotld</code></h1> 
  <p><strong><em>Enumerate all possible root domains of any organization on the fly ✨</em></strong></p>
</div>

# Installation 📩

- Using go ( Assuming you have set `$GOPATH`)
```bash
$ go get -v github.com/WHOISshuvam/gotld
```
- From source
```bash
$ git clone --depth=1 https://github.com/WHOISshuvam/gotld 
$ cd gotld
$ go build . 
```

# Usages 💡
- Helpmenu ( No arguments )

![](/extras/help.png)

- default scans ( No custom keywords)
```
$ gotld -k <keyword> -o <output>
```


![](/extras/withoutfile.png)

- User defined custom keywords file
```
$ gotld -k <keyword> -f <wordlist>  -o <output>
```

![](/extras/withfile.png)

# Thanks 🌺
Thanks to [pwnwriter](https://github.com/pwnwriter) for code/ui `improvement` and  custom `wordlist` options.
