require 'weigh/flags'
require 'weigh/runner'
require 'weigh/util'

module Weigh
  class CLI

    attr_reader :flags
    attr_reader :runner

    def initialize(flags)
      @flags  = flags
      @runner = Weigh::Runner.new(@flags)
    end

    def run
      if flags.help?
        puts flags
        exit
      end

      runner.run
    end

    def self.shutdown
      puts "Terminating..."
      exit 0
    end

    Signal.trap("TERM") do
      shutdown()
    end

    Signal.trap("INT") do
      shutdown()
    end

    def self.run(*args)
      flags = Weigh::Flags.new args
      data = Weigh::CLI.new(flags).run
      Weigh::Util.report(data)
    end

  end
end
