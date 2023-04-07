subdirectory holds project-specific code.
This directory ('apps') is the main module, and the apps within are each
 declared to be part of that module

To set up this module call
 % go mod init main

This needs to be edited so that references to i2s2 are included
when github.com/i-sqr-s-sqr/i2s2 is not set up.  The statement

 replace github.com/i-sqr-s-sqr/i2s2 => /Users/nicol/Dropbox/go-envs/sim-em/i2s2

inside of go.mod will do this

then the call
 % go mod tidy

will discover dependencies

To set up i2s2 call
 % go mod init github.com/i-sqr-s-sqr/i2s2

and then call
 % go mod tidy

Modules are the roots of file directory trees that contain packages.  Intermediate nodes do not need to be packages.

So for instance directory i2s2/multires_network/mrp contains files (one actually) related to package 'mrp',
and each has in its pre-amble 'package mrp', and an application that uses that package imports

 github.com/i-sqr-s-sqr/is2s/multires_network/mrp

so here 'multires_network' just just an intermediate node in the tree of packages rooted in module i2s2



