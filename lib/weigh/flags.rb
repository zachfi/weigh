require 'optparse'

module Weigh
  class Flags

    attr_reader :args
    attr_reader :verbose
    attr_reader :depth
    attr_reader :pathlist

    def initialize(*args)
      @args     = []
      @verbose  = false
      @depth    = 1
      @pathlist = []

      @options = OptionParser.new do|o|
        o.banner = "Usage: #{File.basename $0} [options] [file|directory...]\n\n"

        o.on( '--verbose', '-v', 'Speak up' ) do
          @verbose = true
        end

        o.on( '--depth DEPTH', '-d', 'Sumarize deeper than DEPTH' ) do|d|
          @depth = d
        end

        o.on( '-h', '--help', 'Display this screen' ) do
          @help = true
        end
      end

      @args = @options.parse!

      if ARGV.size > 0
        ARGV.each do |f|
          @pathlist << f
        end
      else
        @pathlist << "."
      end

      def help?
        @help
      end
    end
  end
end
