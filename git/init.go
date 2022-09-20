package git

/*
func InitStageCommitPushToOrigin(
	ctx context.Context,
	repo Local,
	repoURL string,
	stage files.FormFiles,
	commitMsg string,
) error {
	if err := repo.Init(ctx); err != nil {
		return err
	}
	if err := repo.Dir().WriteFormFiles(stage); err != nil {
		return err
	}
	if err := repo.Add(ctx, stage.Paths()...); err != nil {
		return err
	}
	if err := repo.Commit(ctx, commitMsg); err != nil {
		return err
	}
	if err := repo.RenameBranch(ctx, config.MainBranch); err != nil {
		return err
	}
	if err := repo.AddRemoteOrigin(ctx, repoURL); err != nil {
		return err
	}
	return repo.PushToOrigin(ctx, config.MainBranch)
}
*/
