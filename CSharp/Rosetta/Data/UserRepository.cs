using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Query.SqlExpressions;

namespace Rosetta.Data;

public class UserRepository : IUserRepository
{
    private readonly RosettaDbContext _context;
    private readonly ILogger<UserRepository> _logger;
    private readonly TimeSpan _timeout;

    public UserRepository(RosettaDbContext context, ILogger<UserRepository> logger, TimeSpan timeout)
    {
        _context = context;
        _logger = logger;
        _timeout = timeout;
    }

    public async Task<(List<User>?, Metadata?)> SelectAllUsersAsync(Filters filters)
    {
        List<User> users;
        var metadata = new Metadata();
        using var ctx = new CancellationTokenSource(delay: _timeout);
        try
        {
            _logger.LogInformation("Selecting user with filters: {Filters}", filters);
            users = await _context.Users
                .FromSqlInterpolated($@"
SELECT id, name, username, email, created_at, updated_at, deleted, deleted_at
FROM forum.users
WHERE ({filters.ID?.ToString() ?? "NULL"}::UUID IS NULL OR id = {filters.ID}::UUID)
 AND ({filters.Name ?? "NULL"}::VARCHAR(256) IS NULL or name = {filters.Name}::VARCHAR(256))
 AND ({filters.Username ?? "NULL"}::VARCHAR(256) IS NULL or username = {filters.Username}::VARCHAR(256))
 AND ({filters.Email ?? "NULL"}::VARCHAR(256) IS NULL or email = {filters.Email}::VARCHAR(256))
 AND ({filters.CreatedAtFrom?.ToString("yyyy-MM-dd HH:mm:ss") ?? "NULL"}::TIMESTAMP IS NULL or created_at >= {filters.CreatedAtFrom}::TIMESTAMP)
 AND ({filters.CreatedAtTo?.ToString("yyyy-MM-dd HH:mm:ss") ?? "NULL"}::TIMESTAMP IS NULL or created_at <= {filters.CreatedAtTo}::TIMESTAMP)
 AND ({filters.UpdatedAtFrom?.ToString("yyyy-MM-dd HH:mm:ss") ?? "NULL"}::TIMESTAMP IS NULL or updated_at >= {filters.UpdatedAtFrom}::TIMESTAMP)
 AND ({filters.UpdatedAtTo?.ToString("yyyy-MM-dd HH:mm:ss") ?? "NULL"}::TIMESTAMP IS NULL or updated_at <= {filters.UpdatedAtTo}::TIMESTAMP)
{CreateOrderByClause(filters.OrderBy)}
LIMIT {filters.PageSize}::INTEGER;
"
                )
                .AsNoTracking()
                .ToListAsync(cancellationToken: ctx.Token);
        }
        catch (OperationCanceledException)
        {
            _logger.LogError("Query timeout exceeded");
            throw;
        }
        catch (Exception e)
        {
            _logger.LogError(e, "Error performing query");
            throw;
        }

        metadata.ResponseLength = users.Count();
        metadata.Next = users.Count > 0;
        metadata.LastSeen = users.LastOrDefault()?.Id;
        
        _logger.LogInformation("Users selected: {Metadata}", metadata);
        
        return (users, metadata);
    }

    public async Task<User?> SelectUserAsync(Guid id)
    {
        using var ctx = new CancellationTokenSource(delay: _timeout);
        try
        {
            _logger.LogInformation("Selecting user with id: {id}", id);
            var user = await _context.Users
                .FromSqlInterpolated($"""
                    SELECT id, name, username, email, created_at, updated_at, deleted, deleted_at
                    FROM forum.users
                    WHERE id = {id};
                    """
                    )
                .AsNoTracking()
                .FirstOrDefaultAsync(cancellationToken: ctx.Token);
            if (user == null)
            {
                _logger.LogInformation(message: "User with ID {Id} not found", args: id);
                return null;
            }

            _logger.LogInformation(message: "User selected: {User}", args: user);

            return user;
        }
        catch (OperationCanceledException)
        {
            _logger.LogError(message: "Query timeout exceeded for ID: {Id}", args: id);
            throw;
        }
        catch (Exception e)
        {
            _logger.LogError(exception: e, message: "Error performing query: {e}", e);
            throw;
        }
    }

    public async Task<User> UpdateUserAsync(User user)
    {
        _context.Entry(user).State = EntityState.Modified;
        await _context.SaveChangesAsync();
        return user;
    }

    public async Task<User> SoftDeleteUserAsync(Guid id)
    {
        return await _context.Users.FindAsync(id);
    }

    public async Task<User> RestoreUserAsync(Guid id)
    {
        return await _context.Users.FindAsync(id);
    }

    public async Task<User> DeleteUserAsync(Guid id)
    {
        return await _context.Users.FindAsync(id);
    }

    private string CreateOrderByClause(IList<string> orderBy)
    {
        if (!orderBy.Any()) { return "ORDER BY id"; }
        
        IList<string> clauses = new List<string>();
        foreach (var orderByClause in orderBy)
        {
            if (!orderByClause.StartsWith("-"))
            {
                clauses.Add(orderByClause.TrimStart('-') + " DESC");
            }
            else
            {
                clauses.Add(orderByClause + " ASC");
            }
        }
        clauses.Add("id ASC");
        
        return "ORDER BY " + string.Join(", ", clauses);
    }
}