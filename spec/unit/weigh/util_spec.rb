require 'weigh/cli'
require 'weigh/util'
require 'weigh/version'

describe 'Weigh::Util' do
  it "should return a neat size if bytes" do
    Weigh::Util.neat_size(1000).should eq("1000 bytes")
  end

  it "should return a neat size if Kilobytes" do
    Weigh::Util.neat_size(129238).should eq("126.21 KiB")
  end

  it "should return a neat size if Megabytes" do
    Weigh::Util.neat_size(42927238).should eq("40.94 MiB")
  end

  it "should return a neat size if Gigabytes" do
    Weigh::Util.neat_size(85642927238).should eq("79.76 GiB")
  end

  it "should return a neat size if Terabytes" do
    Weigh::Util.neat_size(10285642927238).should eq("9.35 TiB")
  end

  it "should return a neat size if Petabytes" do
    Weigh::Util.neat_size(129210285642927238).should eq("114.76 PiB")
  end
end

describe 'Weigh::CLI' do
  it "should return the version" do
    lambda { Weigh::CLI.command.run(['-V']) }.should raise_error SystemExit
  end
end
