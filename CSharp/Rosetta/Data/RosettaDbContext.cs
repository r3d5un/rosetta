using Microsoft.EntityFrameworkCore;

namespace Rosetta.Data;

public class RosettaDbContext(DbContextOptions<RosettaDbContext> options) : DbContext(options)
{
    // TODO: Add DbSets
}