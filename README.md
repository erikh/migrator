# migrator: migrate SQL databases with the power of... SQL

[![Build Status](https://travis-ci.org/erikh/migrator.svg?branch=master)](https://travis-ci.org/erikh/migrator)<Paste>

migrator takes a different approach to many schema migration systems that exist
already: it hands you the gun. migrator applies SQL directly to your database
in a given order determined by the integers in the filename (minus the `.sql` 
extension) in a provided directory. It is formed from years of experience and
lessons learned writing and using these dumb tools.

migrator has no notion of a "down" migration. This probably isn't changing. If
you wish to alter the tables to an older version of the schema, we strongly
suggest writing a future migration which will transform your data back, because
that is a hell of a lot smarter.

## Example

For example, your database must **already be created**, and your dir should
look like this:

```
dir/
  0.sql
  1.sql
  2.sql
```

The `.sql` files contain arbitrary SQL, but is intended to result in a DDL
transaction where tables are modified, indexes etc. You can use as advanced SQL
as you like and as many statements. Nothing is done with set results, so make
sure they apply to where you want them to!

When finished migrating this directory, your id would be `2`. This will carry
over to the next run to avoid applying those migrations.

To apply new migrations, simply run the run the migrator against the same
directory with newer migrations:

```
dir/
  0.sql
  1.sql
  2.sql
  3.sql
  4.sql
```

In this case, if you already migrated up to `2` you will migrate next to `3`
and `4`. Your schema ID will be `4` after this.

### More Examples

Try:

* Download the repository and make sure docker is installed
* Run `make shell`. https://github.com/box-builder/box will be installed, an
  image will be built and a container will be launched for that image with your
  `$GOPATH` mounted.
* `createdb one` at the shell. This will create a database in postgres called
  `one`.
* `migrator testdata/one` will apply all the migrations for our first test into
  the `one` database.
* If you want different db, or user, etc; check out the flags! `migrator
  --help`.

## License

Copyright 2018 Erik Hollensbe

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Author

Erik Hollensbe <h-e@hollensbe.org>
