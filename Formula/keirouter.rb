# Auto-updated by release.yml on tag v0.1.28. Do not edit manually.
class Keirouter < Formula
  desc "AI API router — unified gateway for 20+ LLM providers with fallback, caching, and dashboard"
  homepage "https://github.com/mydisha/keirouter"
  version "0.1.28"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/mydisha/keirouter/releases/download/v0.1.28/keirouter_v0.1.28_darwin_arm64.tar.gz"
      sha256 "060fba546dab2e7729abdeff54ac26f15275789a4e8f9996442bd5ee3cf305ef"
    else
      url "https://github.com/mydisha/keirouter/releases/download/v0.1.28/keirouter_v0.1.28_darwin_amd64.tar.gz"
      sha256 "c542a406a4a0561d90045fd4f622711c31979c51f06835fe2224b8f2564bb7bc"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/mydisha/keirouter/releases/download/v0.1.28/keirouter_v0.1.28_linux_arm64.tar.gz"
      sha256 "20f41935cd18f9777672e623e0df2352a6ca5f93dbd358c8411faecd95ca435d"
    else
      url "https://github.com/mydisha/keirouter/releases/download/v0.1.28/keirouter_v0.1.28_linux_amd64.tar.gz"
      sha256 "fe4811074603de802ced7cc9cdc0149a767bb0ac72fad671b936a1958d66f0f6"
    end
  end

  def install
    bin.install "keirouter"
    (share/"keirouter").install "frontend"
  end

  def caveats
    <<~EOS
      Quick start:
        keirouter -bootstrap    # create your first API key
        keirouter start         # start server on :20180

      Dashboard: http://localhost:20180  (default password: keirouter)
    EOS
  end

  test do
    assert_match "KeiRouter", shell_output("\#{bin}/keirouter --help 2>&1")
  end
end
