require 'weigh/run'
require 'weigh/util'
require 'cri'

module Weigh::CLI

  def self.command
    @cmd ||= Cri::Command.define do
      name 'weigh'
      usage 'weigh [options] <path1> <path2>'
      summary 'Summarize the size of directories and files'

      w = Weigh::Run.new

      flag :h, :help, 'show help for this command' do |value, cmd|
        puts cmd.help
        exit 0
      end

      flag :V, :version, 'show help for this command' do |value, cmd|
        require 'weigh/version'
        puts "Weigh.rb version " + Weigh::VERSION
        exit 0
      end

      flag :v, :verbose, 'speak up' do |value, cmd|
        w.verbose = true
      end

      flag :d, :depth, 'Sumarize deeper than depth' do |value, cmd|
        #w.depth = value
      end

      run do |opts, args, cmd|
        if args.size > 0
          w.pathlist = args
        end
        w.run
        w.report
        exit 0
      end
    end
  end
end
