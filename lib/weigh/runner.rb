require 'weigh/util'
require 'pp'
require 'find'

module Weigh
  class Runner

    attr_reader :flags

    def initialize(flags)
      @flags = flags
    end

    def run

      # Initialize the data blob
      data              = {}
      data[:total_size] = 0
      data[:summary]    = {}
      data[:count]      = 0

      # Dig into each path given
      @flags.pathlist.each do |p|
        Find.find(p) do |path|
          data[:count] += 1

          # Skip symlinks
          next if FileTest.symlink?(path)

          if FileTest.directory?(path)
            next if p == path
            ret = Weigh::Util.sum_dir(path,@flags.verbose)
            dir_size = ret[:dir_size]

            # Skip empty directories
            next if dir_size == 0

            data[:count]      += ret[:count]
            data[:total_size] += dir_size

            # Add a trailing slash to the key for directories
            pathname = path + "/"

            data[:summary][pathname] = dir_size

            Find.prune
          else
            # Store the size of the current file
            size = FileTest.size(path)

            # Don't count zero sized files
            next if size == 0

            data[:total_size] += size
            data[:summary][path] = size
          end
        end
      end

      data
    end
  end
end
