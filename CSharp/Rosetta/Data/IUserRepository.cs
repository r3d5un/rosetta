namespace Rosetta.Data;

public interface IUserRepository
{
    Task<(List<User>?, Metadata?)> SelectAllUsersAsync(Filters filters);
    Task<User?> SelectUserAsync(Guid id);
    Task<User> UpdateUserAsync(User user);
    Task<User> SoftDeleteUserAsync(Guid id);
    Task<User> RestoreUserAsync(Guid id);
    Task<User> DeleteUserAsync(Guid id);
}