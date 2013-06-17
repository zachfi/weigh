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
      data = {}
      data[:total_size] = 0
      data[:summary] = {}
      data[:count] = 0

      sumdep = @flags.depth

      @flags.pathlist.each do |p|
        curdep = 0
        Find.find(p) do |path|
          data[:count] += 1
          next if FileTest.symlink?(path)
          if FileTest.directory?(path)
            next if p == path
            ret = Weigh::Util.sum_dir(path,@flags.verbose)
            dir_size = ret[:dir_size]
            next if dir_size == 0
            data[:count] += ret[:count]
            data[:total_size] += dir_size
            data[:summary]["#{path}/"] = dir_size
            Find.prune
          else
            size = FileTest.size(path)
            next if size == 0
            data[:total_size] += size
            data[:summary]["#{path}"] = size
          end
        end
      end

      data
    end
  end
end
