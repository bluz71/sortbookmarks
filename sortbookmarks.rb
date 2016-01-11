#!/usr/local/bin/ruby

BOOKMARK_HEADER='<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3 ADD_DATE="1431142608" LAST_MODIFIED="1431142616" PERSONAL_TOOLBAR_FOLDER="true">Bookmarks Bar</H3>
    <DL><p>'
BOOKMARK_FOOTER='    </DL><p>
</DL><p>'

FOLDER_RE = /^\s{8}<DT><H3 ADD_DATE="\d*" LAST_MODIFIED="\d*">(\w+)<\/H3>/
SUBFOLDER_RE = /^\s{12}<DT><H3 ADD_DATE="\d*" LAST_MODIFIED="\d*">(\w+)<\/H3>/
PAGE_RE = /\s+<DT><A HREF="([\w:?&=\-\/\.]+)" ADD_DATE="\d+".*>(.+)<\/A>/

$bookmark_format
$bookmarks = []


class Page
  attr_reader :name, :url

  def initialize(name, url)
    @name = name
    @url = url
  end

  def <=>(page)
    self.name.downcase <=> page.name.downcase
  end
end

class Folder
  attr_reader :name, :subfolders, :pages

  def initialize(name)
    @name = name
    @subfolders = []
    @pages = []
  end

  def <=>(folder)
    self.name.downcase <=> folder.name.downcase
  end
end


def read_and_process_bookmarks
  folder = nil
  subfolder = nil
  add_to_subfolder = false

  File.open("bookmarks.html") do |f|
    f.each do |line|
      matches = FOLDER_RE.match(line)
      if matches
        folder = Folder.new(matches[1])
        $bookmarks << folder
        add_to_subfolder = false
      end
      matches = SUBFOLDER_RE.match(line)
      if matches
        subfolder = Folder.new(matches[1])
        folder.subfolders << subfolder
        add_to_subfolder = true
      end
      matches = PAGE_RE.match(line)
      if matches
        page = Page.new(matches[2], matches[1])
        add_to_subfolder ? subfolder.pages << page : folder.pages << page
      end
    end
  end
end


def print_mozilla_bookmarks
  puts BOOKMARK_HEADER
  $bookmarks.sort.each do |folder|
    puts "        <DT><H3 ADD_DATE=\"1414814291\" LAST_MODIFIED=\"1432192423\">#{folder.name}</H3>"
    puts "        <DL><p>"
    folder.subfolders.sort.each do |subfolder|
      puts "            <DT><H3 ADD_DATE=\"1414814291\" LAST_MODIFIED=\"1432192423\">#{subfolder.name}</H3>"
      puts "            <DL><p>"
      subfolder.pages.sort.each do |page|
        puts "                <DT><A HREF=\"#{page.url}\" ADD_DATE=\"1251002296\">#{page.name}</A>"
      end
      puts "            </DL><p>"
    end
    folder.pages.sort.each do |page|
      puts "            <DT><A HREF=\"#{page.url}\" ADD_DATE=\"1251002296\">#{page.name}</A>"
    end
    puts "        </DL><p>"
  end
  puts BOOKMARK_FOOTER
end


unless ARGV.size == 1
  STDERR.puts "Usage: sortbookmarks.rb mozilla"
  exit
end

unless ARGV[0] == "mozilla"
  STDERR.puts "Unknown sortbookmarks.rb argument: #{ARGV[0]}, expecting 'mozilla'."
  exit
end

unless FileTest.exist?("bookmarks.html")
  STDERR.puts "Expecting 'bookmarks.html' in current directory."
  exit
end

$bookmark_format = ARGV[0]

read_and_process_bookmarks
print_mozilla_bookmarks
