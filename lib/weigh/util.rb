module Weigh
  module Util

    def self.neat_size(bytes)
      # return a human readble size from the bytes supplied
      bytes = bytes.to_f
      if bytes > 2 ** 40       # TiB: 1024 GiB
        neat = sprintf("%.2f TiB", bytes / 2**40)
      elsif bytes > 2 ** 30    # GiB: 1024 MiB
        neat = sprintf("%.2f GiB", bytes / 2**30)
      elsif  bytes > 2 ** 20   # MiB: 1024 KiB
        neat = sprintf("%.2f MiB", bytes / 2**20)
      elsif bytes > 2 ** 10    # KiB: 1024 B
        neat = sprintf("%.2f KiB", bytes / 2**10)
      else                    # bytes
        neat = sprintf("%.0f bytes", bytes)
      end
      neat
    end

  end
end
