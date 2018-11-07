#!/usr/bin/env perl
use 5.014;    # implies "use strict;"
use warnings;
use autodie;

my %total;
my %count;

while (<>) {
    chomp;
    my ( $date, $term, $count ) = split /;/;
    my $month = $date;
    $month =~ s/\-\d+$//;
    $total{$month}{$term} += $count;
    $count{$month}{$term}++;
}

foreach my $month ( sort { $a cmp $b } keys %total ) {
    foreach my $term ( keys %{ $total{$month} } ) {
        my $average = $total{$month}{$term} / $count{$month}{$term};
        printf "%s-01;%s;%.0f\n", $month, $term, $average;
    }
}

__DATA__
2018-11-05;gcp;2
2018-11-05;agile;219
2018-11-05;scrum;165
2018-11-05;sre;438
2018-11-05;mysql;104
2018-11-05;postgresql;37
2018-11-05;oracle;227
2018-11-05;couchdb;1
2018-11-05;redis;17
2018-11-05;mongodb;26
2018-11-04;gcp;2
2018-11-04;agile;219
2018-11-04;scrum;165
2018-11-04;sre;438
2018-11-04;mysql;104
2018-11-04;postgresql;37
2018-11-04;oracle;227
2018-11-04;couchdb;1
2018-11-04;redis;17
2018-11-04;mongodb;26
