pkgname=swaync-widgets-git
pkgver=25.g5352701
pkgrel=1
pkgdesc="A tool for dynamically updating swaync config files based on states"
arch=('x86_64')
url="https://github.com/luiz734/swaync-widgets"
license=('MIT')  # Adjust based on your license
depends=('go')
makedepends=('git' 'go')
provides=($pkgname)
conflicts=($pkgname)
source=("$pkgname::git+https://github.com/luiz734/swaync-widgets")
sha256sums=('SKIP')  # Git sources don't require checksums

pkgver() {
    cd "$srcdir/$pkgname"
    echo "$(git rev-list --count HEAD).g$(git rev-parse --short HEAD)"
}

build() {
    cd "$srcdir/$pkgname"
    export GOPATH="$srcdir"
    go build -o "$srcdir/$pkgname/$pkgname"
}

package() {
    cd "$srcdir/$pkgname"
    install -Dm755 "$pkgname" "$pkgdir/usr/bin/swaync-widgets"
}

