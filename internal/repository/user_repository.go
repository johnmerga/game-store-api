package repository

import (
    "context"
    "database/sql"
    "marketplace-backend/internal/db"
    "marketplace-backend/internal/models"
    
    "github.com/google/uuid"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    Update(ctx context.Context, user *models.User) (*models.User, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error
    List(ctx context.Context, role *models.UserRole, status *models.UserStatus, limit, offset int) ([]*models.User, error)
}

type userRepository struct {
    queries *db.Queries
}

func NewUserRepository(queries *db.Queries) UserRepository {
    return &userRepository{queries: queries}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
    dbUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Email:        user.Email,
        PasswordHash: user.PasswordHash,
        FirstName:    user.FirstName,
        LastName:     user.LastName,
        Role:         db.UserRole(user.Role),
        Phone:        user.Phone,
    })
    if err != nil {
        return nil, err
    }
    
    return r.dbUserToModel(dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    dbUser, err := r.queries.GetUserByEmail(ctx, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    
    return r.dbUserToModel(dbUser), nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    dbUser, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    
    return r.dbUserToModel(dbUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
    dbUser, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
        ID:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Phone:     user.Phone,
        AvatarUrl: user.AvatarURL,
    })
    if err != nil {
        return nil, err
    }
    
    return r.dbUserToModel(dbUser), nil
}

func (r *userRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error {
    return r.queries.UpdateUserStatus(ctx, db.UpdateUserStatusParams{
        ID:     id,
        Status: db.UserStatus(status),
    })
}

func (r *userRepository) List(ctx context.Context, role *models.UserRole, status *models.UserStatus, limit, offset int) ([]*models.User, error) {
    var dbRole *db.UserRole
    var dbStatus *db.UserStatus
    
    if role != nil {
        dbRoleVal := db.UserRole(*role)
        dbRole = &dbRoleVal
    }
    if status != nil {
        dbStatusVal := db.UserStatus(*status)
        dbStatus = &dbStatusVal
    }
    
    dbUsers, err := r.queries.ListUsers(ctx, db.ListUsersParams{
        Role:   dbRole,
        Status: dbStatus,
