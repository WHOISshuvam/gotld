<div align="center">
  <h1><code>gotld</code></h1> 
  <p><strong><em>Enumerate all possible root domains of any organization on the fly ✨</em></strong></p>
</div>

# Installation 📩

- Using go ( Assuming you have set `$GOPATH`)
```bash
$ go install github.com/WHOISshuvam/gotld@latest
```
- From source
```bash
$ git clone --depth=1 https://github.com/WHOISshuvam/gotld 
$ cd gotld
$ go build . 
```

# Usages 💡
- Help Menu ( No arguments )

![](/extras/help.png)

- default scans ( No custom keywords)
```
$ gotld -k <keyword> -o <output>
```


![](/extras/withoutfile.png)

- User defined custom `wordlist` file
```
$ gotld -k <keyword> -f <wordlist>  -o <output>
```

![](/extras/withfile.png)

- Number of threads to use `(default 5)`
```
$ gotld -k <keyword> -t <number of threads> -o <output>
```
![](/extras/threads.png)

# Thanks 🌺
Thanks to [pwnwriter](https://github.com/pwnwriter) for code/ui `improvement` and  custom `wordlist` options.
