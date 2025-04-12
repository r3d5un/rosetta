using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace Rosetta.Data;

[Table("users", Schema = "forum")]
public class User
{
    [
        Column("id", TypeName = "UUID"),
        Required,
    ]
    public Guid Id { get; set; }
    
    [
        Column("name", TypeName = "VARCHAR(256)"),
        Required,
    ]
    public string Name { get; set; }
    
    [
        Column("username", TypeName = "VARCHAR(256)"),
        Required,
    ]
    public string Username { get; set; }
    
    [
        Column("email", TypeName = "VARCHAR(256)"),
        Required,
    ]
    public string Email { get; set; }
    
    [
        Column("created_at", TypeName = "TIMESTAMP"),
        Required,
    ]
    public DateTime CreatedAt { get; set; }
    
    [
        Column("updated_at", TypeName = "TIMESTAMP"),
        Required,
    ]
    public DateTime UpdatedAt { get; set; }
    
    [
        Column("deleted", TypeName = "BOOLEAN"),
        Required,
    ]
    public bool Deleted { get; set; }
    
    [
    Column("deleted_at", TypeName = "TIMESTAMP"),
    Required,
    ]
    public DateTime DeletedAt { get; set; }
    
    public List<User> Users { get; set; } = [];
}