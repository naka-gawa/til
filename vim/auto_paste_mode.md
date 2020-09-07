# vimで自動でpaste modeにする

## Overview
タイトルどおり。vimのコピペ時にインデントが崩れるので
`:set paste`を都度実行していた。
でも続いて編集する場合は自動インデントは必要なので
更に`:set nopaste`を実行したりして、つらみがあった。

## 暫定対応
```
:set noautoindent
```

```
:set paste
```

```
:a!
```

上記はどれもpaste modeにすることができるが、手動での実行になるので後者の辛みは解消しない。

## Neobundle
```
cd ~/.vim/bundle
git clone https://github.com/ConradIrwin/vim-bracketed-paste

or

NeoBundle 'ConradIrwin/vim-bracketed-paste'
```
もしくは、.vimrcに直接
```
if &term =~ "xterm"
    let &t_ti .= "\e[?2004h"
    let &t_te .= "\e[?2004l"
    let &pastetoggle = "\e[201~"

    function XTermPasteBegin(ret)
        set paste
        return a:ret
    endfunction

    noremap <special> <expr> <Esc>[200~ XTermPasteBegin("0i")
    inoremap <special> <expr> <Esc>[200~ XTermPasteBegin("")
    cnoremap <special> <Esc>[200~ <nop>
    cnoremap <special> <Esc>[201~ <nop>
endif
```
で解消する

### 補足事項
- bracketed paste modeに対応しているターミナルが必要

>You need to be using a modern xterm-compatible terminal emulator that supports bracketed paste mode. xterm, urxvt, iTerm2, gnome-terminal (and other terminals using libvte) are known to work.

## ref
[vimでペーストする際に、自動でpaste modeにする方法のメモ](https://qiita.com/ryoff/items/ad34584e41425362453e)
