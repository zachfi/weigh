module Weigh
  class Runner
    attr_reader :flags

    def initialize(flags)
      @flags = flags
    end

    def run
      puts "no"
      return 127
    end

  end
end
