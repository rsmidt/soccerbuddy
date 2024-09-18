package domain

type Repositories interface {
	Account() AccountRepository
	Club() ClubRepository
	Person() PersonRepository
	Session() SessionRepository
	Team() TeamRepository
	TeamMember() TeamMemberRepository
}
