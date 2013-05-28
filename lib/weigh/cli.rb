require 'weigh/flags'
require 'weigh/runner'

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

    def self.run(*args)
      flags = Weigh::Flags.new args

      return Weigh::CLI.new(flags).run
    end
  end
end
