# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Envexample < Formula
  desc "Generate a .env.example from an env struct"
  homepage "https://github.com/miniscruff/envexample"
  version "0.1.2"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/miniscruff/envexample/releases/download/v0.1.2/envexample_0.1.2_darwin_arm64.tar.gz"
      sha256 "c6377a3a23bb6894b917fa2d006f1ea42d898618cd282deb41246af38966f5d2"

      def install
        bin.install "envexample"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/miniscruff/envexample/releases/download/v0.1.2/envexample_0.1.2_darwin_amd64.tar.gz"
      sha256 "68f6ab42616076b970fd36d475df1dd71eba153850a32ab79ccf2168c5ac2dac"

      def install
        bin.install "envexample"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/miniscruff/envexample/releases/download/v0.1.2/envexample_0.1.2_linux_amd64.tar.gz"
      sha256 "b2ec3c90d5ec8d85ee97dac087adcac12f43473b6d066059165a25721059643a"

      def install
        bin.install "envexample"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/miniscruff/envexample/releases/download/v0.1.2/envexample_0.1.2_linux_arm64.tar.gz"
      sha256 "8889584ad75c1579c44b4662d2808982814c2f9a2542ff3888be7cad5d7c31fd"

      def install
        bin.install "envexample"
      end
    end
  end
end
