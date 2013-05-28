require 'weigh/report'

module Weigh
  class Runner

    attr_reader :flags

    def initialize(flags)
      @flags = flags
    end

    def run
      summary = {}
      Weigh::Report.print(summary)

      return 127
    end

  end
end
