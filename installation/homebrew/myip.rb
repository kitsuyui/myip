require "formula"

class Myip < Formula
  homepage "https://github.com/kitsuyui/go-myip"

  if Hardware::CPU.is_64_bit?
    url "https://github.com/kitsuyui/go-myip/releases/download/v0.2.1/myip_darwin_amd64"
    sha256 "f28d142dd8063c789e8539e4ede21c87a23ef4b91614a1294d2c102888407181"
  else
    url "https://github.com/kitsuyui/go-myip/releases/download/v0.2.1/myip_darwin_386"
    sha256 "d431dba6ec931abee0c65f90c9535274a01388b0e7bc32e445ba5fa3ab650354"
  end

  head "https://github.com/kitsuyui/go-myip.git"
  version "v0.2.1"

  def install
    if Hardware::CPU.is_64_bit?
      bin.install "myip_darwin_amd64" => "myip"
    else
      bin.install "myip_darwin_386" => "myip"
    end
  end
end
