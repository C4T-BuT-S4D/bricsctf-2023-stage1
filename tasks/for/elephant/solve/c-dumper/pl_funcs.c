/*-------------------------------------------------------------------------
 *
 * pl_funcs.c		- Misc functions for the PL/pgSQL
 *			  procedural language
 *
 * Portions Copyright (c) 1996-2023, PostgreSQL Global Development Group
 * Portions Copyright (c) 1994, Regents of the University of California
 *
 *
 * IDENTIFICATION
 *	  src/pl/plpgsql/src/pl_funcs.c
 *
 *-------------------------------------------------------------------------
 */
#include "stdio.h"

#include "plpgsql_ida.h"
#include "pg_list.h"

#define _

/*
 * Statement type as a string, for use in error messages etc.
 */
const char *
plpgsql_stmt_typename(PLpgSQL_stmt *stmt)
{
	switch (stmt->cmd_type)
	{
		case PLPGSQL_STMT_BLOCK:
			return _("statement block");
		case PLPGSQL_STMT_ASSIGN:
			return _("assignment");
		case PLPGSQL_STMT_IF:
			return "IF";
		case PLPGSQL_STMT_CASE:
			return "CASE";
		case PLPGSQL_STMT_LOOP:
			return "LOOP";
		case PLPGSQL_STMT_WHILE:
			return "WHILE";
		case PLPGSQL_STMT_FORI:
			return _("FOR with integer loop variable");
		case PLPGSQL_STMT_FORS:
			return _("FOR over SELECT rows");
		case PLPGSQL_STMT_FORC:
			return _("FOR over cursor");
		case PLPGSQL_STMT_FOREACH_A:
			return _("FOREACH over array");
		case PLPGSQL_STMT_EXIT:
			return ((PLpgSQL_stmt_exit *) stmt)->is_exit ? "EXIT" : "CONTINUE";
		case PLPGSQL_STMT_RETURN:
			return "RETURN";
		case PLPGSQL_STMT_RETURN_NEXT:
			return "RETURN NEXT";
		case PLPGSQL_STMT_RETURN_QUERY:
			return "RETURN QUERY";
		case PLPGSQL_STMT_RAISE:
			return "RAISE";
		case PLPGSQL_STMT_ASSERT:
			return "ASSERT";
		case PLPGSQL_STMT_EXECSQL:
			return _("SQL statement");
		case PLPGSQL_STMT_DYNEXECUTE:
			return "EXECUTE";
		case PLPGSQL_STMT_DYNFORS:
			return _("FOR over EXECUTE statement");
		case PLPGSQL_STMT_GETDIAG:
			return ((PLpgSQL_stmt_getdiag *) stmt)->is_stacked ?
				"GET STACKED DIAGNOSTICS" : "GET DIAGNOSTICS";
		case PLPGSQL_STMT_OPEN:
			return "OPEN";
		case PLPGSQL_STMT_FETCH:
			return ((PLpgSQL_stmt_fetch *) stmt)->is_move ? "MOVE" : "FETCH";
		case PLPGSQL_STMT_CLOSE:
			return "CLOSE";
		case PLPGSQL_STMT_PERFORM:
			return "PERFORM";
		case PLPGSQL_STMT_CALL:
			return ((PLpgSQL_stmt_call *) stmt)->is_call ? "CALL" : "DO";
		case PLPGSQL_STMT_COMMIT:
			return "COMMIT";
		case PLPGSQL_STMT_ROLLBACK:
			return "ROLLBACK";
	}

	return "unknown";
}

/*
 * GET DIAGNOSTICS item name as a string, for use in error messages etc.
 */
const char *
plpgsql_getdiag_kindname(PLpgSQL_getdiag_kind kind)
{
	switch (kind)
	{
		case PLPGSQL_GETDIAG_ROW_COUNT:
			return "ROW_COUNT";
		case PLPGSQL_GETDIAG_ROUTINE_OID:
			return "PG_ROUTINE_OID";
		case PLPGSQL_GETDIAG_CONTEXT:
			return "PG_CONTEXT";
		case PLPGSQL_GETDIAG_ERROR_CONTEXT:
			return "PG_EXCEPTION_CONTEXT";
		case PLPGSQL_GETDIAG_ERROR_DETAIL:
			return "PG_EXCEPTION_DETAIL";
		case PLPGSQL_GETDIAG_ERROR_HINT:
			return "PG_EXCEPTION_HINT";
		case PLPGSQL_GETDIAG_RETURNED_SQLSTATE:
			return "RETURNED_SQLSTATE";
		case PLPGSQL_GETDIAG_COLUMN_NAME:
			return "COLUMN_NAME";
		case PLPGSQL_GETDIAG_CONSTRAINT_NAME:
			return "CONSTRAINT_NAME";
		case PLPGSQL_GETDIAG_DATATYPE_NAME:
			return "PG_DATATYPE_NAME";
		case PLPGSQL_GETDIAG_MESSAGE_TEXT:
			return "MESSAGE_TEXT";
		case PLPGSQL_GETDIAG_TABLE_NAME:
			return "TABLE_NAME";
		case PLPGSQL_GETDIAG_SCHEMA_NAME:
			return "SCHEMA_NAME";
	}

	return "unknown";
}


/**********************************************************************
 * Debug functions for analyzing the compiled code
 **********************************************************************/
static int	dump_indent;

static void dump_ind(void);
static void dump_stmt(PLpgSQL_stmt *stmt);
static void dump_block(PLpgSQL_stmt_block *block);
static void dump_assign(PLpgSQL_stmt_assign *stmt);
static void dump_if(PLpgSQL_stmt_if *stmt);
static void dump_case(PLpgSQL_stmt_case *stmt);
static void dump_loop(PLpgSQL_stmt_loop *stmt);
static void dump_while(PLpgSQL_stmt_while *stmt);
static void dump_fori(PLpgSQL_stmt_fori *stmt);
static void dump_fors(PLpgSQL_stmt_fors *stmt);
static void dump_forc(PLpgSQL_stmt_forc *stmt);
static void dump_foreach_a(PLpgSQL_stmt_foreach_a *stmt);
static void dump_exit(PLpgSQL_stmt_exit *stmt);
static void dump_return(PLpgSQL_stmt_return *stmt);
static void dump_return_next(PLpgSQL_stmt_return_next *stmt);
static void dump_return_query(PLpgSQL_stmt_return_query *stmt);
static void dump_raise(PLpgSQL_stmt_raise *stmt);
static void dump_assert(PLpgSQL_stmt_assert *stmt);
static void dump_execsql(PLpgSQL_stmt_execsql *stmt);
static void dump_dynexecute(PLpgSQL_stmt_dynexecute *stmt);
static void dump_dynfors(PLpgSQL_stmt_dynfors *stmt);
static void dump_getdiag(PLpgSQL_stmt_getdiag *stmt);
static void dump_open(PLpgSQL_stmt_open *stmt);
static void dump_fetch(PLpgSQL_stmt_fetch *stmt);
static void dump_cursor_direction(PLpgSQL_stmt_fetch *stmt);
static void dump_close(PLpgSQL_stmt_close *stmt);
static void dump_perform(PLpgSQL_stmt_perform *stmt);
static void dump_call(PLpgSQL_stmt_call *stmt);
static void dump_commit(PLpgSQL_stmt_commit *stmt);
static void dump_rollback(PLpgSQL_stmt_rollback *stmt);
static void dump_expr(PLpgSQL_expr *expr);


static void
dump_ind(void)
{
	int			i;

	for (i = 0; i < dump_indent; i++)
		printf(" ");
}

static void
dump_stmt(PLpgSQL_stmt *stmt)
{
	printf("%3d:", stmt->lineno);
	switch (stmt->cmd_type)
	{
		case PLPGSQL_STMT_BLOCK:
			dump_block((PLpgSQL_stmt_block *) stmt);
			break;
		case PLPGSQL_STMT_ASSIGN:
			dump_assign((PLpgSQL_stmt_assign *) stmt);
			break;
		case PLPGSQL_STMT_IF:
			dump_if((PLpgSQL_stmt_if *) stmt);
			break;
		case PLPGSQL_STMT_CASE:
			dump_case((PLpgSQL_stmt_case *) stmt);
			break;
		case PLPGSQL_STMT_LOOP:
			dump_loop((PLpgSQL_stmt_loop *) stmt);
			break;
		case PLPGSQL_STMT_WHILE:
			dump_while((PLpgSQL_stmt_while *) stmt);
			break;
		case PLPGSQL_STMT_FORI:
			dump_fori((PLpgSQL_stmt_fori *) stmt);
			break;
		case PLPGSQL_STMT_FORS:
			dump_fors((PLpgSQL_stmt_fors *) stmt);
			break;
		case PLPGSQL_STMT_FORC:
			dump_forc((PLpgSQL_stmt_forc *) stmt);
			break;
		case PLPGSQL_STMT_FOREACH_A:
			dump_foreach_a((PLpgSQL_stmt_foreach_a *) stmt);
			break;
		case PLPGSQL_STMT_EXIT:
			dump_exit((PLpgSQL_stmt_exit *) stmt);
			break;
		case PLPGSQL_STMT_RETURN:
			dump_return((PLpgSQL_stmt_return *) stmt);
			break;
		case PLPGSQL_STMT_RETURN_NEXT:
			dump_return_next((PLpgSQL_stmt_return_next *) stmt);
			break;
		case PLPGSQL_STMT_RETURN_QUERY:
			dump_return_query((PLpgSQL_stmt_return_query *) stmt);
			break;
		case PLPGSQL_STMT_RAISE:
			dump_raise((PLpgSQL_stmt_raise *) stmt);
			break;
		case PLPGSQL_STMT_ASSERT:
			dump_assert((PLpgSQL_stmt_assert *) stmt);
			break;
		case PLPGSQL_STMT_EXECSQL:
			dump_execsql((PLpgSQL_stmt_execsql *) stmt);
			break;
		case PLPGSQL_STMT_DYNEXECUTE:
			dump_dynexecute((PLpgSQL_stmt_dynexecute *) stmt);
			break;
		case PLPGSQL_STMT_DYNFORS:
			dump_dynfors((PLpgSQL_stmt_dynfors *) stmt);
			break;
		case PLPGSQL_STMT_GETDIAG:
			dump_getdiag((PLpgSQL_stmt_getdiag *) stmt);
			break;
		case PLPGSQL_STMT_OPEN:
			dump_open((PLpgSQL_stmt_open *) stmt);
			break;
		case PLPGSQL_STMT_FETCH:
			dump_fetch((PLpgSQL_stmt_fetch *) stmt);
			break;
		case PLPGSQL_STMT_CLOSE:
			dump_close((PLpgSQL_stmt_close *) stmt);
			break;
		case PLPGSQL_STMT_PERFORM:
			dump_perform((PLpgSQL_stmt_perform *) stmt);
			break;
		case PLPGSQL_STMT_CALL:
			dump_call((PLpgSQL_stmt_call *) stmt);
			break;
		case PLPGSQL_STMT_COMMIT:
			dump_commit((PLpgSQL_stmt_commit *) stmt);
			break;
		case PLPGSQL_STMT_ROLLBACK:
			dump_rollback((PLpgSQL_stmt_rollback *) stmt);
			break;
		default:
      *(void*) 0;
			break;
	}
}

#define false (char) 0
#define true (char) 1
#define NULL 0

static void
dump_stmts(List *stmts)
{
	ListCell   *s;

	dump_indent += 2;
	foreach(s, stmts) {
		dump_stmt((PLpgSQL_stmt *) lfirst(s));
  }
	dump_indent -= 2;
}

static void
dump_block(PLpgSQL_stmt_block *block)
{
	char	   *name;

	if (block->label == NULL)
		name = "*unnamed*";
	else
		name = block->label;

	dump_ind();
	printf("BLOCK <<%s>>\n", name);

	dump_stmts(block->body);

	if (block->exceptions)
	{
		ListCell   *e;

		foreach(e, block->exceptions->exc_list)
		{
			PLpgSQL_exception *exc = (PLpgSQL_exception *) lfirst(e);
			PLpgSQL_condition *cond;

			dump_ind();
			printf("    EXCEPTION WHEN ");
			for (cond = exc->conditions; cond; cond = cond->next)
			{
				if (cond != exc->conditions)
					printf(" OR ");
				printf("%s", cond->condname);
			}
			printf(" THEN\n");
			dump_stmts(exc->action);
		}
	}

	dump_ind();
	printf("    END -- %s\n", name);
}

static void
dump_assign(PLpgSQL_stmt_assign *stmt)
{
	dump_ind();
	printf("ASSIGN var %d := ", stmt->varno);
	dump_expr(stmt->expr);
	printf("\n");
}

static void
dump_if(PLpgSQL_stmt_if *stmt)
{
	ListCell   *l;

	dump_ind();
	printf("IF ");
	dump_expr(stmt->cond);
	printf(" THEN\n");
	dump_stmts(stmt->then_body);
	foreach(l, stmt->elsif_list)
	{
		PLpgSQL_if_elsif *elif = (PLpgSQL_if_elsif *) lfirst(l);

		dump_ind();
		printf("    ELSIF ");
		dump_expr(elif->cond);
		printf(" THEN\n");
		dump_stmts(elif->stmts);
	}
	if (stmt->else_body != NIL)
	{
		dump_ind();
		printf("    ELSE\n");
		dump_stmts(stmt->else_body);
	}
	dump_ind();
	printf("    ENDIF\n");
}

static void
dump_case(PLpgSQL_stmt_case *stmt)
{
	ListCell   *l;

	dump_ind();
	printf("CASE %d ", stmt->t_varno);
	if (stmt->t_expr)
		dump_expr(stmt->t_expr);
	printf("\n");
	dump_indent += 6;
	foreach(l, stmt->case_when_list)
	{
		PLpgSQL_case_when *cwt = (PLpgSQL_case_when *) lfirst(l);

		dump_ind();
		printf("WHEN ");
		dump_expr(cwt->expr);
		printf("\n");
		dump_ind();
		printf("THEN\n");
		dump_indent += 2;
		dump_stmts(cwt->stmts);
		dump_indent -= 2;
	}
	if (stmt->have_else)
	{
		dump_ind();
		printf("ELSE\n");
		dump_indent += 2;
		dump_stmts(stmt->else_stmts);
		dump_indent -= 2;
	}
	dump_indent -= 6;
	dump_ind();
	printf("    ENDCASE\n");
}

static void
dump_loop(PLpgSQL_stmt_loop *stmt)
{
	dump_ind();
	printf("LOOP\n");

	dump_stmts(stmt->body);

	dump_ind();
	printf("    ENDLOOP\n");
}

static void
dump_while(PLpgSQL_stmt_while *stmt)
{
	dump_ind();
	printf("WHILE ");
	dump_expr(stmt->cond);
	printf("\n");

	dump_stmts(stmt->body);

	dump_ind();
	printf("    ENDWHILE\n");
}

static void
dump_fori(PLpgSQL_stmt_fori *stmt)
{
	dump_ind();
	printf("FORI %s %s\n", stmt->var->refname, (stmt->reverse) ? "REVERSE" : "NORMAL");

	dump_indent += 2;
	dump_ind();
	printf("    lower = ");
	dump_expr(stmt->lower);
	printf("\n");
	dump_ind();
	printf("    upper = ");
	dump_expr(stmt->upper);
	printf("\n");
	if (stmt->step)
	{
		dump_ind();
		printf("    step = ");
		dump_expr(stmt->step);
		printf("\n");
	}
	dump_indent -= 2;

	dump_stmts(stmt->body);

	dump_ind();
	printf("    ENDFORI\n");
}

static void
dump_fors(PLpgSQL_stmt_fors *stmt)
{
	dump_ind();
	printf("FORS %s ", stmt->var->refname);
	dump_expr(stmt->query);
	printf("\n");

	dump_stmts(stmt->body);

	dump_ind();
	printf("    ENDFORS\n");
}

static void
dump_forc(PLpgSQL_stmt_forc *stmt)
{
	dump_ind();
	printf("FORC %s ", stmt->var->refname);
	printf("curvar=%d\n", stmt->curvar);

	dump_indent += 2;
	if (stmt->argquery != NULL)
	{
		dump_ind();
		printf("  arguments = ");
		dump_expr(stmt->argquery);
		printf("\n");
	}
	dump_indent -= 2;

	dump_stmts(stmt->body);

	dump_ind();
	printf("    ENDFORC\n");
}

static void
dump_foreach_a(PLpgSQL_stmt_foreach_a *stmt)
{
	dump_ind();
	printf("FOREACHA var %d ", stmt->varno);
	if (stmt->slice != 0)
		printf("SLICE %d ", stmt->slice);
	printf("IN ");
	dump_expr(stmt->expr);
	printf("\n");

	dump_stmts(stmt->body);

	dump_ind();
	printf("    ENDFOREACHA");
}

static void
dump_open(PLpgSQL_stmt_open *stmt)
{
	dump_ind();
	printf("OPEN curvar=%d\n", stmt->curvar);

	dump_indent += 2;
	if (stmt->argquery != NULL)
	{
		dump_ind();
		printf("  arguments = '");
		dump_expr(stmt->argquery);
		printf("'\n");
	}
	if (stmt->query != NULL)
	{
		dump_ind();
		printf("  query = '");
		dump_expr(stmt->query);
		printf("'\n");
	}
	if (stmt->dynquery != NULL)
	{
		dump_ind();
		printf("  execute = '");
		dump_expr(stmt->dynquery);
		printf("'\n");

		if (stmt->params != NIL)
		{
			ListCell   *lc;
			int			i;

			dump_indent += 2;
			dump_ind();
			printf("    USING\n");
			dump_indent += 2;
			i = 1;
			foreach(lc, stmt->params)
			{
				dump_ind();
				printf("    parameter $%d: ", i++);
				dump_expr((PLpgSQL_expr *) lfirst(lc));
				printf("\n");
			}
			dump_indent -= 4;
		}
	}
	dump_indent -= 2;
}

static void
dump_fetch(PLpgSQL_stmt_fetch *stmt)
{
	dump_ind();

	if (!stmt->is_move)
	{
		printf("FETCH curvar=%d\n", stmt->curvar);
		dump_cursor_direction(stmt);

		dump_indent += 2;
		if (stmt->target != NULL)
		{
			dump_ind();
			printf("    target = %d %s\n",
				   stmt->target->dno, stmt->target->refname);
		}
		dump_indent -= 2;
	}
	else
	{
		printf("MOVE curvar=%d\n", stmt->curvar);
		dump_cursor_direction(stmt);
	}
}

static void
dump_cursor_direction(PLpgSQL_stmt_fetch *stmt)
{
	dump_indent += 2;
	dump_ind();
	switch (stmt->direction)
	{
		case FETCH_FORWARD:
			printf("    FORWARD ");
			break;
		case FETCH_BACKWARD:
			printf("    BACKWARD ");
			break;
		case FETCH_ABSOLUTE:
			printf("    ABSOLUTE ");
			break;
		case FETCH_RELATIVE:
			printf("    RELATIVE ");
			break;
		default:
			printf("??? unknown cursor direction %d", stmt->direction);
	}

	if (stmt->expr)
	{
		dump_expr(stmt->expr);
		printf("\n");
	}
	else
		printf("%ld\n", stmt->how_many);

	dump_indent -= 2;
}

static void
dump_close(PLpgSQL_stmt_close *stmt)
{
	dump_ind();
	printf("CLOSE curvar=%d\n", stmt->curvar);
}

static void
dump_perform(PLpgSQL_stmt_perform *stmt)
{
	dump_ind();
	printf("PERFORM expr = ");
	dump_expr(stmt->expr);
	printf("\n");
}

static void
dump_call(PLpgSQL_stmt_call *stmt)
{
	dump_ind();
	printf("%s expr = ", stmt->is_call ? "CALL" : "DO");
	dump_expr(stmt->expr);
	printf("\n");
}

static void
dump_commit(PLpgSQL_stmt_commit *stmt)
{
	dump_ind();
	if (stmt->chain)
		printf("COMMIT AND CHAIN\n");
	else
		printf("COMMIT\n");
}

static void
dump_rollback(PLpgSQL_stmt_rollback *stmt)
{
	dump_ind();
	if (stmt->chain)
		printf("ROLLBACK AND CHAIN\n");
	else
		printf("ROLLBACK\n");
}

static void
dump_exit(PLpgSQL_stmt_exit *stmt)
{
	dump_ind();
	printf("%s", stmt->is_exit ? "EXIT" : "CONTINUE");
	if (stmt->label != NULL)
		printf(" label='%s'", stmt->label);
	if (stmt->cond != NULL)
	{
		printf(" WHEN ");
		dump_expr(stmt->cond);
	}
	printf("\n");
}

static void
dump_return(PLpgSQL_stmt_return *stmt)
{
	dump_ind();
	printf("RETURN ");
	if (stmt->retvarno >= 0)
		printf("variable %d", stmt->retvarno);
	else if (stmt->expr != NULL)
		dump_expr(stmt->expr);
	else
		printf("NULL");
	printf("\n");
}

static void
dump_return_next(PLpgSQL_stmt_return_next *stmt)
{
	dump_ind();
	printf("RETURN NEXT ");
	if (stmt->retvarno >= 0)
		printf("variable %d", stmt->retvarno);
	else if (stmt->expr != NULL)
		dump_expr(stmt->expr);
	else
		printf("NULL");
	printf("\n");
}

static void
dump_return_query(PLpgSQL_stmt_return_query *stmt)
{
	dump_ind();
	if (stmt->query)
	{
		printf("RETURN QUERY ");
		dump_expr(stmt->query);
		printf("\n");
	}
	else
	{
		printf("RETURN QUERY EXECUTE ");
		dump_expr(stmt->dynquery);
		printf("\n");
		if (stmt->params != NIL)
		{
			ListCell   *lc;
			int			i;

			dump_indent += 2;
			dump_ind();
			printf("    USING\n");
			dump_indent += 2;
			i = 1;
			foreach(lc, stmt->params)
			{
				dump_ind();
				printf("    parameter $%d: ", i++);
				dump_expr((PLpgSQL_expr *) lfirst(lc));
				printf("\n");
			}
			dump_indent -= 4;
		}
	}
}

static void
dump_raise(PLpgSQL_stmt_raise *stmt)
{
	ListCell   *lc;
	int			i = 0;

	dump_ind();
	printf("RAISE level=%d", stmt->elog_level);
	if (stmt->condname)
		printf(" condname='%s'", stmt->condname);
	if (stmt->message)
		printf(" message='%s'", stmt->message);
	printf("\n");
	dump_indent += 2;
	foreach(lc, stmt->params)
	{
		dump_ind();
		printf("    parameter %d: ", i++);
		dump_expr((PLpgSQL_expr *) lfirst(lc));
		printf("\n");
	}
	if (stmt->options)
	{
		dump_ind();
		printf("    USING\n");
		dump_indent += 2;
		foreach(lc, stmt->options)
		{
			PLpgSQL_raise_option *opt = (PLpgSQL_raise_option *) lfirst(lc);

			dump_ind();
			switch (opt->opt_type)
			{
				case PLPGSQL_RAISEOPTION_ERRCODE:
					printf("    ERRCODE = ");
					break;
				case PLPGSQL_RAISEOPTION_MESSAGE:
					printf("    MESSAGE = ");
					break;
				case PLPGSQL_RAISEOPTION_DETAIL:
					printf("    DETAIL = ");
					break;
				case PLPGSQL_RAISEOPTION_HINT:
					printf("    HINT = ");
					break;
				case PLPGSQL_RAISEOPTION_COLUMN:
					printf("    COLUMN = ");
					break;
				case PLPGSQL_RAISEOPTION_CONSTRAINT:
					printf("    CONSTRAINT = ");
					break;
				case PLPGSQL_RAISEOPTION_DATATYPE:
					printf("    DATATYPE = ");
					break;
				case PLPGSQL_RAISEOPTION_TABLE:
					printf("    TABLE = ");
					break;
				case PLPGSQL_RAISEOPTION_SCHEMA:
					printf("    SCHEMA = ");
					break;
			}
			dump_expr(opt->expr);
			printf("\n");
		}
		dump_indent -= 2;
	}
	dump_indent -= 2;
}

static void
dump_assert(PLpgSQL_stmt_assert *stmt)
{
	dump_ind();
	printf("ASSERT ");
	dump_expr(stmt->cond);
	printf("\n");

	dump_indent += 2;
	if (stmt->message != NULL)
	{
		dump_ind();
		printf("    MESSAGE = ");
		dump_expr(stmt->message);
		printf("\n");
	}
	dump_indent -= 2;
}

static void
dump_execsql(PLpgSQL_stmt_execsql *stmt)
{
	dump_ind();
	printf("EXECSQL ");
	dump_expr(stmt->sqlstmt);
	printf("\n");

	dump_indent += 2;
	if (stmt->target != NULL)
	{
		dump_ind();
		printf("    INTO%s target = %d %s\n",
			   stmt->strict ? " STRICT" : "",
			   stmt->target->dno, stmt->target->refname);
	}
	dump_indent -= 2;
}

static void
dump_dynexecute(PLpgSQL_stmt_dynexecute *stmt)
{
	dump_ind();
	printf("EXECUTE ");
	dump_expr(stmt->query);
	printf("\n");

	dump_indent += 2;
	if (stmt->target != NULL)
	{
		dump_ind();
		printf("    INTO%s target = %d %s\n",
			   stmt->strict ? " STRICT" : "",
			   stmt->target->dno, stmt->target->refname);
	}
	if (stmt->params != NIL)
	{
		ListCell   *lc;
		int			i;

		dump_ind();
		printf("    USING\n");
		dump_indent += 2;
		i = 1;
		foreach(lc, stmt->params)
		{
			dump_ind();
			printf("    parameter %d: ", i++);
			dump_expr((PLpgSQL_expr *) lfirst(lc));
			printf("\n");
		}
		dump_indent -= 2;
	}
	dump_indent -= 2;
}

static void
dump_dynfors(PLpgSQL_stmt_dynfors *stmt)
{
	dump_ind();
	printf("FORS %s EXECUTE ", stmt->var->refname);
	dump_expr(stmt->query);
	printf("\n");
	if (stmt->params != NIL)
	{
		ListCell   *lc;
		int			i;

		dump_indent += 2;
		dump_ind();
		printf("    USING\n");
		dump_indent += 2;
		i = 1;
		foreach(lc, stmt->params)
		{
			dump_ind();
			printf("    parameter $%d: ", i++);
			dump_expr((PLpgSQL_expr *) lfirst(lc));
			printf("\n");
		}
		dump_indent -= 4;
	}
	dump_stmts(stmt->body);
	dump_ind();
	printf("    ENDFORS\n");
}

static void
dump_getdiag(PLpgSQL_stmt_getdiag *stmt)
{
	ListCell   *lc;

	dump_ind();
	printf("GET %s DIAGNOSTICS ", stmt->is_stacked ? "STACKED" : "CURRENT");
	foreach(lc, stmt->diag_items)
	{
		PLpgSQL_diag_item *diag_item = (PLpgSQL_diag_item *) lfirst(lc);

		if (lc != list_head(stmt->diag_items))
			printf(", ");

		printf("{var %d} = %s", diag_item->target,
			   plpgsql_getdiag_kindname(diag_item->kind));
	}
	printf("\n");
}

static void
dump_expr(PLpgSQL_expr *expr)
{
	printf("'%s'", expr->query);
}

void
plpgsql_dumptree(PLpgSQL_function *func)
{
	int			i;
	PLpgSQL_datum *d;

	printf("\nExecution tree of successfully compiled PL/pgSQL function %s:\n",
		   func->fn_signature);

	printf("\nFunction's data area:\n");
	for (i = 0; i < func->ndatums; i++)
	{
		d = func->datums[i];

		printf("    entry %d: ", i);
		switch (d->dtype)
		{
			case PLPGSQL_DTYPE_VAR:
			case PLPGSQL_DTYPE_PROMISE:
				{
					PLpgSQL_var *var = (PLpgSQL_var *) d;

					printf("VAR %-16s type %s (typoid %u) atttypmod %d\n",
						   var->refname, var->datatype->typname,
						   var->datatype->typoid,
						   var->datatype->atttypmod);
					if (var->isconst)
						printf("                                  CONSTANT\n");
					if (var->notnull)
						printf("                                  NOT NULL\n");
					if (var->default_val != NULL)
					{
						printf("                                  DEFAULT ");
						dump_expr(var->default_val);
						printf("\n");
					}
					if (var->cursor_explicit_expr != NULL)
					{
						if (var->cursor_explicit_argrow >= 0)
							printf("                                  CURSOR argument row %d\n", var->cursor_explicit_argrow);

						printf("                                  CURSOR IS ");
						dump_expr(var->cursor_explicit_expr);
						printf("\n");
					}
					if (var->promise != PLPGSQL_PROMISE_NONE)
						printf("                                  PROMISE %d\n",
							   (int) var->promise);
				}
				break;
			case PLPGSQL_DTYPE_ROW:
				{
					PLpgSQL_row *row = (PLpgSQL_row *) d;

					printf("ROW %-16s fields", row->refname);
					for (int j = 0; j < row->nfields; j++)
					{
						printf(" %s=var %d", row->fieldnames[j],
							   row->varnos[j]);
					}
					printf("\n");
				}
				break;
			case PLPGSQL_DTYPE_REC:
				printf("REC %-16s typoid %u\n",
					   ((PLpgSQL_rec *) d)->refname,
					   ((PLpgSQL_rec *) d)->rectypeid);
				if (((PLpgSQL_rec *) d)->isconst)
					printf("                                  CONSTANT\n");
				if (((PLpgSQL_rec *) d)->notnull)
					printf("                                  NOT NULL\n");
				if (((PLpgSQL_rec *) d)->default_val != NULL)
				{
					printf("                                  DEFAULT ");
					dump_expr(((PLpgSQL_rec *) d)->default_val);
					printf("\n");
				}
				break;
			case PLPGSQL_DTYPE_RECFIELD:
				printf("RECFIELD %-16s of REC %d\n",
					   ((PLpgSQL_recfield *) d)->fieldname,
					   ((PLpgSQL_recfield *) d)->recparentno);
				break;
			default:
				printf("??? unknown data type %d\n", d->dtype);
		}
	}
	printf("\nFunction's statements:\n");

	dump_indent = 0;
	printf("%3d:", func->action->lineno);
	dump_block(func->action);
	printf("\nEnd of execution tree of function %s\n\n", func->fn_signature);
	fflush(stdout);
}
